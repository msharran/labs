set -x

for i in $(seq 1 10); do
    # out=$(curl -s --resolve foo.example.com:80:192.168.49.2 foo.example.com)
    out=$(curl -s -H "Host: foo.example.com" http://172.18.0.2:30587)
    if [[ $out == *"404"* ]]; then
            echo "path not found"
    elif [[ $out == *"nginx"* ]]; then
        echo "nginx"
    elif [[ $out == *"works!"* ]]; then
        echo "httpd (canary)"
    fi
done
