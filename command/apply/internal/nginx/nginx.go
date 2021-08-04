package nginx

import "fmt"

var conf = `server {
    listen      80;
    listen      [::]:80;
    set         $base /app;
    root        $base/%s;

    proxy_send_timeout 240s;
    proxy_read_timeout 240s;
    fastcgi_send_timeout 240s;
    fastcgi_read_timeout 240s;

    # security
    include     nitro/security.conf;

    # include custom conf files
    include     /app/*nitro.conf;

    # index.php
    index       index.php;

    # index.php fallback
    location / {
        try_files $uri $uri/ /index.php?$query_string;
    }

    # additional config
    include nitro/general.conf;

    # handle .php
    location ~ \.php$ {
        include nitro/php_fastcgi.conf;
    }

    # Allow fpm ping and status from localhost
    location ~ ^/(fpm-status|fpm-ping)$ {
        access_log off;
        allow 127.0.0.1;
        deny all;
        fastcgi_param SCRIPT_FILENAME $document_root$fastcgi_script_name;
        include fastcgi_params;
        fastcgi_pass 127.0.0.1:9000;
    }
}`

// Generate takes a root directory generates a nginx configuration file
func Generate(root string) string {
	// if the root was not provided, default to web
	if root == "" {
		return fmt.Sprintf(conf, "web")
	}

	return fmt.Sprintf(conf, root)
}
