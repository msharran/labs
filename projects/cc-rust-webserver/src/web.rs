pub mod http;
use http::{HttpRequest, HttpResponse};
use std::{collections::HashMap, fs, net::TcpStream};

pub struct TcpConnManager<'a> {
    stream: &'a mut TcpStream,
    html_files: HashMap<String, String>, // name -> content
}

impl<'a> TcpConnManager<'a> {
    pub fn from(stream: &'a mut TcpStream) -> Self {
        let mut html_files = HashMap::new();
        let index_content =
            fs::read_to_string("www/index.html").expect("Unable to read index.html");
        html_files.insert("index.html".to_string(), index_content);

        Self { stream, html_files }
    }

    pub fn handle_connection(&mut self) {
        println!("HTTP Stream accepted");

        let http_request = match HttpRequest::from(&mut self.stream) {
            Ok(req) => req,
            Err(e) => {
                eprintln!("Error: {e}");
                return;
            }
        };

        println!("{http_request:?}");

        let response = match http_request.uri.as_str() {
            "/" | "/index.html" => {
                let content = self.html_files.get(&"index.html".to_string());
                if let Some(content) = content {
                    HttpResponse::ok(content.to_string())
                } else {
                    HttpResponse::not_found("Page not found".to_string())
                }
            }
            _ => HttpResponse::not_found("Page not found".to_string()),
        };

        response.write_all(self.stream);
        println!("HTTP Stream closed");
    }
}
