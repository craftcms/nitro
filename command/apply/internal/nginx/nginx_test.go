package nginx

import "testing"

func TestGenerate(t *testing.T) {
	type args struct {
		root string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "defaults to web",
			args: args{
				root: "",
			},
			want: defaultConf,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Generate(tt.args.root); got != tt.want {
				t.Errorf("Generate() = %v, want %v", got, tt.want)
			}
		})
	}
}

var defaultConf = `server {
    listen      8080 default_server;
    listen      [::]:8080 default_server;
    server_name _;
    set         $base /app;
    root        $base/web;

    proxy_send_timeout 240s;
    proxy_read_timeout 240s;
    fastcgi_send_timeout 240s;
    fastcgi_read_timeout 240s;

    # security
    include     craftcms/security.conf;

    # include custom conf files
    include     /app/*nitro.conf;

    # index.php
    index       index.php;

    # index.php fallback
    location / {
        try_files $uri $uri/ /index.php?$query_string;
    }

    # additional config
    include craftcms/general.conf;

    # handle .php
    location ~ \.php$ {
        include craftcms/php_fastcgi.conf;
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
