set -x

curl -H "Host: jenkins-agent.example.com" http://172.18.0.2:30587/ | w3m -dump -T text/html
