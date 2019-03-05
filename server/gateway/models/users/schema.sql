create table IF NOT EXISTS users (
  userID int not null auto_increment primary key,
  email varchar(256) UNIQUE not null ,
  username varchar(256) UNIQUE not null,
  passHash BINARY(60) not null,
  firstName varchar(50) not null,
  lastName varchar(50) not null,
  photoUrl varchar(2083) not NULL
);

create table IF NOT EXISTS channels (
  channelID int not null AUTO_INCREMENT primary key,
  channelName varchar(256) UNIQUE not null,
  description varchar(512),
  privateChannel BOOLEAN not null,
  createdAt DATETIME not null,
  creatorID int not null,
  editedAt DATETIME
);

create table IF NOT EXISTS channel_user (
  channelUserID int not null AUTO_INCREMENT primary key,
  userID int not null,
  channelID int not NULL
);

create table IF NOT EXISTS messages(
  messageID int not null AUTO_INCREMENT primary key,
  channelID int not NULL,
  body varchar(2048),
  createdAt DATETIME not null,
  creatorID int not null,
  editedAt DATETIME
);

create table IF NOT EXISTS starred_message (
  starID int not null AUTO_INCREMENT primary key,
  messageID int not NULL,
  userID int not NULL,
  channelID int not NULL
);

INSERT INTO channels(channelName, description, privateChannel, createdAt, creatorID) VALUES("General", "This is a general channel for everyone", false, CURRENT_TIMESTAMP(),0);
