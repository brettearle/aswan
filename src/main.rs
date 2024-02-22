use std::{
    env,
    io::{self, Read},
};

fn new_project(args: &Vec<String>) {
    //check args for project name and git link, ensure they are present
    //if not, print error message and exit
    if args.len() < 4 || args.len() > 4 {
        println!("Usage: aswan <new/ls> <project_name> <git_link>");
        return;
    }
    let input = Project::new(&args);
    //write project to local file
    match std::fs::write(
        "projects.txt",
        format!("{}, {}\n", input.project_name, input.git_link),
    ) {
        Ok(file) => file,
        Err(e) => panic!("Error writing to file: {}", e),
    };
    //open project.txt
    std::process::Command::new("code")
        .arg("projects.txt")
        .spawn()
        .expect("Failed to open file");
}

fn list_projects() {
    let file = match std::fs::read_to_string("projects.txt") {
        Ok(file) => file,
        Err(e) => panic!("Error reading file: {}", e),
    };
    println!("{}", file);
}

fn main() {
    let args: Vec<String> = env::args().collect();
    let mut command = String::new();

    io::stdin()
        .read_line(&mut command)
        .expect("Failed to read line");

    println!("Command: {command}");

    let director = get_director(&args);
    //check args for director command
    //if not present, print error message and exit
    match director {
        DirectorArguments::New => new_project(&args),
        DirectorArguments::Ls => list_projects(),
    }
}

struct Project {
    project_name: String,
    git_link: String,
}

impl Project {
    fn new(args_vec: &Vec<String>) -> Project {
        Project {
            project_name: args_vec[2].clone(),
            git_link: args_vec[3].clone(),
        }
    }
}

enum DirectorArguments {
    New,
    Ls,
}

fn get_director(args: &Vec<String>) -> DirectorArguments {
    match args[1].as_str() {
        "new" => DirectorArguments::New,
        "ls" => DirectorArguments::Ls,
        _ => panic!("Usage: aswan <new/ls>"),
    }
}
