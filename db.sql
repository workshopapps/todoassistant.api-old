CREATE DATABASE IF NOT EXISTS todoDB;

USE todoDB;

DROP TABLES Users, Virtual_Assistants, Tasks, Virtual_Assistant_Ratings, Taskfiles, Calls;

CREATE TABLE IF NOT EXISTS Users(
	`user_id` VARCHAR(255) NOT NULL UNIQUE,
    `first_name` VARCHAR(255) NOT NULL,
    `last_name` VARCHAR(255) NOT NULL,
    `email` VARCHAR(255) NOT NULL UNIQUE,
    `phone` CHAR NOT NULL UNIQUE,
    `password` VARCHAR(255) NOT NULL,
    `gender` ENUM('Male', 'Female', 'Others'),
    `date_of_birth` CHAR NOT NULL,
    `virtual_assistant_id` VARCHAR(255),
    `account_status` ENUM('Active', 'Suspended', 'Blocked'),
    `payment_status` ENUM('Active', 'Suspended', 'Blocked'),
    `date_created` DATETIME default CURRENT_TIMESTAMP,
    `date_updated` DATETIME default CURRENT_TIMESTAMP,
    `last_login` DATETIME default CURRENT_TIMESTAMP,
    
    PRIMARY KEY (user_id)
);

CREATE TABLE IF NOT EXISTS Virtual_Assistants(
	`virtual_assistant_id` VARCHAR(255) NOT NULL UNIQUE,
    `first_name` VARCHAR(255) NOT NULL,
    `last_name` VARCHAR(255) NOT NULL,
    `email` VARCHAR(255) NOT NULL UNIQUE,
    `phone` CHAR NOT NULL UNIQUE,
    `password` VARCHAR(255) NOT NULL,
    `date_of_birth` CHAR NOT NULL,
    `account_status` ENUM('Active', 'Suspended', 'Blocked'),
    `date_created` DATETIME default CURRENT_TIMESTAMP,
    `date_updated` DATETIME default CURRENT_TIMESTAMP,
    `last_login` DATETIME default CURRENT_TIMESTAMP,
    
    PRIMARY KEY (virtual_assistant_id)
);

CREATE TABLE IF NOT EXISTS Tasks(
	`task_id` VARCHAR(255) NOT NULL UNIQUE,
    `user_id` VARCHAR(255) NOT NULL,
	`title` CHAR NOT NULL,
    `description` TEXT NOT NULL,
    `status` ENUM ('Pending', 'Done', 'Expired', 'Overdue') DEFAULT 'Overdue',
	`start_time` DATETIME default CURRENT_TIMESTAMP,
	`end_time` DATETIME NOT NULL,
    `created_at` CHAR(50),
    `updated_at` CHAR(50),
    
    PRIMARY KEY(task_id),
    FOREIGN KEY (user_id) REFERENCES Users(user_id)
);

CREATE TABLE IF NOT EXISTS Virtual_Assistant_Ratings(
	`id` VARCHAR(255) NOT NULL,
    `user_id` VARCHAR(255) NOT NULL,
    `va_id` VARCHAR(255) NOT NULL,
    `ratings` INT NOT NULL default 5,
    `date` DATETIME default CURRENT_TIMESTAMP,
    
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES Users(user_id),
    FOREIGN KEY (va_id) REFERENCES Virtual_Assistants(virtual_assistant_id)
);

CREATE TABLE IF NOT EXISTS Taskfiles(
	`id` INT NOT NULL AUTO_INCREMENT,
    `task_id` VARCHAR(255) NOT NULL,
    `file_link` CHAR NOT NULL,
    `file_type` ENUM('Audio', 'Video', 'Image') NOT NULL,
    
    PRIMARY KEY (id),
    FOREIGN KEY (task_id) REFERENCES Tasks(task_id)
);


CREATE TABLE IF NOT EXISTS Calls(
	`id` VARCHAR(255) NOT NULL,
    `user_id` VARCHAR(255) NOT NULL,
    `va_id` VARCHAR(255) NOT NULL,
	`call_rating` INT,
    `call_comment` CHAR,
    
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES Users(user_id),
    FOREIGN KEY (va_id) REFERENCES Virtual_Assistants(virtual_assistant_id)
);