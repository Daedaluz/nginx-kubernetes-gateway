package config

var serversTemplateText = `
{{ range $s := . -}}
	{{ if $s.IsDefaultSSL -}}
server {
	listen 443 ssl default_server;
    listen [::]:443 ssl default_server;

	ssl_reject_handshake on;
}
	{{- else if $s.IsDefaultHTTP }}
server {
	listen 80 default_server;
    listen [::]:80 default_server;

	default_type text/html;
	return 404;
}
	{{- else }}
server {
		{{- if $s.SSL }}
	listen 443 ssl;
    listen [::]:443 ssl;

	ssl_certificate {{ $s.SSL.Certificate }};
	ssl_certificate_key {{ $s.SSL.CertificateKey }};

	if ($ssl_server_name != $host) {
		return 421;
	}
		{{- end }}

	server_name {{ $s.ServerName }};
	listen 80;
    listen [::]:80;

		{{ range $l := $s.Locations }}
	location {{ if $l.Exact }}= {{ end }}{{ $l.Path }} {
		{{ if $l.Internal -}}
		internal;
		{{ end }}

		{{- if $l.Return -}}
		return {{ $l.Return.Code }} "{{ $l.Return.Body }}";
		{{ end }}

		{{- if $l.HTTPMatchVar -}}
		set $http_matches {{ $l.HTTPMatchVar | printf "%q" }};
		js_content httpmatches.redirect;
		{{ end }}

		{{- if $l.ProxyPass -}}
		proxy_pass {{ $l.ProxyPass }}$request_uri;
		proxy_set_header        Host               $host;
		proxy_set_header        X-Real-IP          $remote_addr;
		proxy_set_header        X-Forwarded-For    $proxy_add_x_forwarded_for;
		proxy_set_header        X-Forwarded-Host   $host:443;
		proxy_set_header        X-Forwarded-Server $host;
		proxy_set_header        X-Forwarded-Port   443;
		proxy_set_header        X-Forwarded-Proto  https;
		{{- end }}
	}
		{{ end }}
}
	{{- end }}
{{ end }}
server {
    listen unix:/var/lib/nginx/nginx-502-server.sock;
    access_log off;

    return 502;
}

server {
    listen unix:/var/lib/nginx/nginx-500-server.sock;
    access_log off;
    
    return 500;
}
`
