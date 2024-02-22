use std::env;

fn new_project(args: &Vec<String>) {
    //check args for project name and git link, ensure they are present
    //if not, print error message and exit
    if args.len() < 3 || args.len() > 3 {
        println!("Usage: aswan <project_name> <git_link>");
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
    let director = get_director(&args);
    //check args for director command
    //if not present, print error message and exit
    match director {
        DirectorArguments::New => new_project(&args),
        DirectorArguments::Ls => list_projects(),
    }
    let input = Project::new(&args);
    //write project to local file
    let file = match std::fs::write(
        "projects.txt",
        format!("{}, {}\n", input.project_name, input.git_link),
    ) {
        Ok(file) => file,
        Err(e) => panic!("Error writing to file: {}", e),
    };

    println!("Project Name: {}", input.project_name);
    println!("Git Link: {}", input.git_link);
}

struct Project {
    project_name: String,
    git_link: String,
}

impl Project {
    fn new(args_vec: &Vec<String>) -> Project {
        Project {
            project_name: args_vec[1].clone(),
            git_link: args_vec[2].clone(),
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
