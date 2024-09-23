use std::{
    io::{BufRead, BufReader, Write},
    net::{TcpListener, TcpStream},
};

fn main() {
    let socket = TcpListener::bind("127.0.0.1:8080").unwrap();
    println!("Socket bound to 8080");

    for stream in socket.incoming() {
        match stream {
            Err(e) => {
                eprintln!("Err: Cannot establish connection, {e:?}");
                continue;
            }
            Ok(s) => {
                handle_http_connection(s);
            }
        };
    }
}

fn handle_http_connection(mut stream: TcpStream) {
    println!("Connection accepted");
    let reader = BufReader::new(&mut stream);

    let http_req_lines = read_http_request(reader);
    match http_req_lines {
        Err(e) => {
            eprintln!("Err: Cannot read http request {e:?}");
            return;
        }
        Ok(lines) => {
            println!("Request: {lines:?}");

            let req = parse_http_request(lines);
            println!("{req:?}");

            let response = "HTTP/1.1 200 OK\r\n\r\n";
            stream.write_all(response.as_bytes()).unwrap();
        }
    };
}

fn read_http_request(r: BufReader<&mut TcpStream>) -> Result<Vec<String>, String> {
    let mut errs = vec![];
    let req_lines: Vec<_> = r
        .lines()
        .filter_map(|r| r.map_err(|e| errs.push(e.to_string())).ok())
        .take_while(|l| !l.is_empty())
        .collect();

    if errs.len() > 0 {
        return Err(errs.join(": "));
    }

    Ok(req_lines)
}

#[derive(Debug)]
struct HttpRequest {
    // headers: HashMap<String, String>,
    version: String,
    uri: String,
    method: String,
}

fn parse_http_request(r: Vec<String>) -> Result<HttpRequest, String> {
    // lines won't be empty at any cost
    // as we are filterning only non empty lines
    let fline = r.iter().next().unwrap();
    let fparts: Vec<_> = fline.split_whitespace().collect();
    if fparts.len() != 3 {
        return Err(format!(
            "request first line must be 3 parts. got: {}",
            fline
        ));
    }

    let method = fparts.get(0).unwrap().to_string();
    let uri = fparts.get(1).unwrap().to_string();
    let version = fparts.get(2).unwrap().to_string();

    let req = HttpRequest {
        uri,
        method,
        version,
    };
    Ok(req)
}
