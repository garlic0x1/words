use std::{io, io::prelude::*};
use url::Url;

fn main() {
    reader();
}

fn reader() {
    for line in io::stdin().lock().lines() {
        parse_line(line.unwrap_or_default());
    }
}

fn parse_line(line: String) {
    match Url::parse(line.as_str()) {
        Ok(u) => {
            let furl = serde_json::json!({
                "scheme": u.scheme(),
                "domain": u.host().unwrap().to_string(),
                "path": u.path(),
                "query_str": u.query().unwrap_or_default(),
            });

            println!("{:}", furl);
        },
        Err(_e) => {
            return;
        },
    }
}
