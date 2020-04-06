package nitro

func AddSiteScript(name, domain, php, dir string) Command {
	return Command{
		Machine: name,
		Type:    "exec",
		Args:    []string{"--", "sudo", "bash", "/opt/nitro/nginx/add-site.sh", domain, php, dir},
	}
}
