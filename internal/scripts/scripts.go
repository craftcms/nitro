package scripts

import (
	"fmt"
	"os/exec"
	"strings"
)

const (
	FmtNginxSiteAvailable                     = `if test -f '/etc/nginx/sites-available/%s'; then echo 'exists'; fi`
	FmtNginxSiteEnabled                       = `if test -f '/etc/nginx/sites-enabled/%s'; then echo 'exists'; fi`
	FmtNginxSiteWebroot                       = `grep "root " /etc/nginx/sites-available/%s | while read -r line; do echo "$line"; done`
	FmtDockerContainerExists                  = `if [ -n "$(docker ps -q -f name="%s")" ]; then echo "exists"; fi`
	FmtDockerMysqlCreateDatabaseIfNotExists   = `docker exec -i %s mysql -unitro -pnitro -e "CREATE DATABASE IF NOT EXISTS %s;"`
	FmtDockerPostgresCreateDatabase           = `docker exec -i %s psql --username nitro -c "CREATE DATABASE %s;"`
	FmtDockerMysqlImportDatabase              = `cat %s | docker exec -i %s mysql -unitro -pnitro %s --init-command="SET autocommit=0;"`
	FmtDockerMysqlShowAllDatabases            = `docker exec -i %s mysql -unitro -pnitro -e "SHOW DATABASES;"`
	FmtDockerPostgresShowAllDatabases         = `docker exec -i %s psql --username nitro --command "SELECT datname FROM pg_database WHERE datistemplate = false;"`
	DockerListContainerNames                  = `docker container ls --all --format '{{ .Names }}'`
	FmtDockerRestartContainer                 = `docker container restart %s`
	FmtDockerStopContainer                    = `docker container stop %s`
	FmtDockerRemoveContainer                  = `docker container rm -f -v %s`
	FmtDockerRemoveVolume                     = `docker volume rm -f %s`
	FmtDockerStartContainer                   = `docker container start %s`
	FmtDockerBackupAllMysqlDatabases          = `docker exec %s /usr/bin/mysqldump --all-databases -unitro -pnitro > %s`
	FmtDockerBackupIndividualPostgresDatabase = `docker exec -i %s pg_dump -U nitro %s > %s`
	FmtDockerBackupIndividualMysqlDatabase    = `docker exec %s /usr/bin/mysqldump -unitro -pnitro %s > %s`
	FmtCreateDirectory                        = `mkdir -p %s`
)

type Script struct {
	path    string
	machine string
}

// New will return a new Script struct that
// contains the path to multipass and
// the name of the machine
func New(multipass, machine string) *Script {
	return &Script{
		path:    multipass,
		machine: machine,
	}
}

// Run is used to make running scripts on a nitro machine
// a lot easier, using New will store the path to the
// nitro path and machine name. Run will then run
// the script on the machine and
func (s Script) Run(sudo bool, arg ...string) (string, error) {
	args := []string{"exec", s.machine, "--"}
	switch sudo {
	case true:
		args = append(args, []string{"sudo", "bash", "-c"}...)
	default:
		args = append(args, []string{"bash", "-c"}...)
	}
	args = append(args, arg...)

	cmd := exec.Command(s.path, args...)

	bytes, err := cmd.CombinedOutput()
	output := strings.TrimSpace(string(bytes))
	if err != nil {
		fmt.Println(output)
		return "", err
	}

	return strings.TrimSpace(output), nil
}
