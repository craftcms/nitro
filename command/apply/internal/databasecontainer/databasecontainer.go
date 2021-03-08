package databasecontainer

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/labels"
	"github.com/craftcms/nitro/pkg/terminal"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	volumetypes "github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	_ "github.com/go-sql-driver/mysql"
)

var (
	// DatabaseImage is used for determining the engine and version
	DatabaseImage = "%s:%s"
)

// StartOrCreate is used to find a specific database and start the container. If there is no container for the database,
// it will create a new volume and container for the database.
func StartOrCreate(ctx context.Context, docker client.CommonAPIClient, networkID string, db config.Database, output terminal.Outputer) (string, string, error) {
	// create the filters for the database
	filter := filters.NewArgs()
	filter.Add("label", labels.DatabaseEngine+"="+db.Engine)
	filter.Add("label", labels.DatabaseVersion+"="+db.Version)
	filter.Add("label", labels.DatabasePort+"="+db.Port)
	filter.Add("label", labels.Type+"=database")

	hostname, err := db.GetHostname()
	if err != nil {
		return "", "", err
	}

	// set the container database compatibility
	if db.Engine == "mariadb" || db.Engine == "mysql" {
		filter.Add("label", labels.DatabaseCompatibility+"=mysql")
	} else {
		filter.Add("label", labels.DatabaseCompatibility+"=postgres")
	}

	// get the containers for the database
	containers, err := docker.ContainerList(ctx, types.ContainerListOptions{All: true, Filters: filter})
	if err != nil {
		return "", "", fmt.Errorf("error getting a list of containers")
	}

	// if there is a container, we should start it and return
	if len(containers) == 1 {
		// check if the container is running
		if containers[0].State != "running" {
			// start the container
			if err := docker.ContainerStart(ctx, containers[0].ID, types.ContainerStartOptions{}); err != nil {
				return "", "", err
			}
		}

		return containers[0].ID, hostname, nil
	}

	// create the database labels for the new container
	lbls := map[string]string{
		labels.Nitro:           "true",
		labels.DatabaseEngine:  db.Engine,
		labels.DatabaseVersion: db.Version,
		labels.Type:            "database",
		labels.DatabasePort:    db.Port,
	}

	// if the database is mysql or mariadb, mark them as
	// mysql compatible (used for importing backups)
	if db.Engine == "mariadb" || db.Engine == "mysql" {
		lbls[labels.DatabaseCompatibility] = "mysql"
	}

	// if the database is postgres, mark it as compatible
	// with postgres. This is not needed but a place holder
	// if cockroachdb is ever supported by craft.
	if db.Engine == "postgres" {
		lbls[labels.DatabaseCompatibility] = "postgres"
	}

	// create the volume
	volume, err := docker.VolumeCreate(ctx, volumetypes.VolumeCreateBody{Driver: "local", Name: hostname, Labels: lbls})
	if err != nil {
		return "", "", fmt.Errorf("unable to create the volume, %w", err)
	}

	// determine the image name
	image := fmt.Sprintf(DatabaseImage, db.Engine, db.Version)

	// set mounts and environment based on the database type
	target := "/var/lib/mysql"
	var envs []string
	if strings.Contains(image, "postgres") {
		target = "/var/lib/postgresql/data"
		envs = []string{"POSTGRES_USER=nitro", "POSTGRES_DB=nitro", "POSTGRES_PASSWORD=nitro"}
	} else {
		envs = []string{"MYSQL_ROOT_PASSWORD=nitro", "MYSQL_DATABASE=nitro", "MYSQL_USER=nitro", "MYSQL_PASSWORD=nitro"}
	}

	// check if there is an image

	// filter for the image ref
	imageFilter := filters.NewArgs()
	imageFilter.Add("reference", image)

	// look for the image
	images, err := docker.ImageList(ctx, types.ImageListOptions{Filters: imageFilter, All: true})
	if err != nil {
		return "", "", fmt.Errorf("unable to get a list of images, %w", err)
	}

	// if there are no images, pull one
	if len(images) == 0 {
		output.Pending("downloading", image)

		// pull the image
		rdr, err := docker.ImagePull(ctx, image, types.ImagePullOptions{All: false})
		if err != nil {
			output.Warning()

			return "", "", fmt.Errorf("unable to pull image %s, %w", image, err)
		}

		// read the output to pull the image
		buf := &bytes.Buffer{}
		if _, err := buf.ReadFrom(rdr); err != nil {
			output.Warning()
			return "", "", fmt.Errorf("unable to read output from pulling image %s, %w", image, err)
		}
	}

	// get the default port for the database
	var port nat.Port
	switch db.Engine {
	case "postgres":
		port, err = nat.NewPort("tcp", "5432")
		if err != nil {
			return "", "", fmt.Errorf("unable to create the port, %w", err)
		}
	default:
		port, err = nat.NewPort("tcp", "3306")
		if err != nil {
			return "", "", fmt.Errorf("unable to create the port, %w", err)
		}
	}

	containerConfig := &container.Config{
		Image:  image,
		Labels: lbls,
		ExposedPorts: nat.PortSet{
			port: struct{}{},
		},
		Env: envs,
	}

	// if the mysql engine is being used, override the cmd
	if db.Engine == "mysql" {
		containerConfig.Cmd = []string{"--character-set-server=utf8mb4", "--collation-server=utf8mb4_unicode_ci"}
	}

	hostConfig := &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeVolume,
				Source: volume.Name,
				Target: target,
			},
		},
		PortBindings: map[nat.Port][]nat.PortBinding{
			port: {
				{
					HostIP:   "127.0.0.1",
					HostPort: db.Port,
				},
			},
		},
	}

	networkConfig := &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			"nitro-network": {
				NetworkID: networkID,
			},
		},
	}

	// create the container for the database
	resp, err := docker.ContainerCreate(ctx, containerConfig, hostConfig, networkConfig, nil, hostname)
	if err != nil {
		return "", "", fmt.Errorf("unable to create the container, %w", err)
	}

	// start the container
	if err := docker.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return "", "", fmt.Errorf("unable to start the container, %w", err)
	}

	// if the container is mysql compatible
	if db.Engine == "mysql" || db.Engine == "mariadb" {
		if err := waitForMySQLContainer(ctx, docker, resp.ID, db); err != nil {
			return "", "", err
		}
	}

	return resp.ID, hostname, nil
}

func waitForMySQLContainer(ctx context.Context, docker client.CommonAPIClient, containerID string, d config.Database) error {
	// verify the mysql socket exists in the container
	for {
		// wait until the socket is ready
		stat, _ := docker.ContainerStatPath(ctx, containerID, "/var/run/mysqld/mysqld.sock")
		if stat.Name != "" {
			break
		}
	}

	// connect to the database
	db, err := sql.Open("mysql", fmt.Sprintf("root:nitro@tcp(127.0.0.1:%s)/nitro", d.Port))
	if err != nil {
		return fmt.Errorf("error opening connection: %w", err)
	}

	// set the connection time to 5 seconds
	db.SetConnMaxLifetime(5 * time.Second)
	db.SetMaxIdleConns(0)
	db.SetMaxOpenConns(151)

	// ugh, sleep for 10 seconds because of mysql...
	wait := time.Duration(time.Second * 10)
	time.Sleep(wait)

	// setup the commands
	commands := [][]string{
		{"mysql", "-uroot", "-pnitro", fmt.Sprintf(`-e CREATE USER IF NOT EXISTS '%s'@'%s' IDENTIFIED BY 'nitro';`, "nitro", "localhost")},
		{"mysql", "-uroot", "-pnitro", fmt.Sprintf(`-e GRANT ALL PRIVILEGES ON *.* TO '%s'@'%s' WITH GRANT OPTION;`, "nitro", "%")},
		{"mysql", "-uroot", "-pnitro", fmt.Sprintf(`-e GRANT ALL PRIVILEGES ON *.* TO '%s'@'%s' WITH GRANT OPTION;`, "nitro", "localhost")},
		{"mysql", "-uroot", "-pnitro", `-e FLUSH PRIVILEGES;`},
	}

	// for mysql 8.0 images
	// ALTER USER ‘username’@‘ip_address’ IDENTIFIED WITH mysql_native_password BY ‘password’
	if strings.Contains(d.Version, "8.0") {
		commands = append(commands, []string{"mysql", "-uroot", "-pnitro", fmt.Sprintf(`-e ALTER USER '%s'@'%s' IDENTIFIED WITH mysql_native_password BY 'nitro';`, "nitro", "%")})
	}

	for _, c := range commands {
		// create the exec
		exec, err := docker.ContainerExecCreate(ctx, containerID, types.ExecConfig{
			AttachStdout: true,
			AttachStderr: true,
			Tty:          false,
			Cmd:          c,
		})
		if err != nil {
			return err
		}

		// attach to the container
		resp, err := docker.ContainerExecAttach(ctx, exec.ID, types.ExecStartCheck{
			Tty: false,
		})
		if err != nil {
			return err
		}

		// start the exec
		if err := docker.ContainerExecStart(ctx, exec.ID, types.ExecStartCheck{}); err != nil {
			return fmt.Errorf("unable to start the container, %w", err)
		}

		// wait for the container exec to complete
		for {
			resp, err := docker.ContainerExecInspect(ctx, exec.ID)
			if err != nil {
				return err
			}

			if !resp.Running {
				break
			}
		}

		// close the exec attach
		resp.Close()
	}

	db.Close()

	return nil
}
