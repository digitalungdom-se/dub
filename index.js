/* global base_dir abs_path client */

global.base_dir = __dirname;
global.abs_path = function ( path ) {
  return `${base_dir}/${path}`;
};
global.include = function ( file ) {
  return require( abs_path( file ) );
};

global.musicQueue = [];
global.musicVolume = 0.1;
global.searchList = false;
global.searchMessage = false;

require( 'dotenv' ).config();
const config = require( './config.json' );

const Discord = require( 'discord.js' );
const Spotify = require( 'node-spotify-api' );
const fs = require( 'fs-extra' );
const MongoClient = require( 'mongodb' ).MongoClient;



// create a new Discord client
global.client = new Discord.Client();
client.commands = new Discord.Collection();

global.spotify = new Spotify( {
  id: process.env.SPOTIPY_CLIENT_ID,
  secret: process.env.SPOTIPY_CLIENT_SECRET
} );

// load message commands
const commandDirectories = fs.readdirSync( './commands' ).filter( function ( file ) {
  return fs.statSync( `./commands/${file}` ).isDirectory();
} );

for ( const commandDirectory of commandDirectories ) {
  const commandFiles = ( fs.readdirSync( `./commands/${commandDirectory}` ).filter( file => file.endsWith( '.js' ) ) );
  for ( const file of commandFiles ) {
    const command = require( `./commands/${commandDirectory}/${file}` );
    client.commands.set( command.name, command );
  }
}

// when the client is ready, run this code
// this event will only trigger one time after logging in
client.once( 'ready', () => {
  console.log( 'Ready!' );
  client.user.setActivity( 'Kelvin\'s cat', { type: 'WATCHING' } );
  global.guild = client.guilds.get( process.env.GUILD_ID );
} );

// mesage functions
client.on( 'message', message => {
  if ( config.prefix.indexOf( message.content.charAt( 0 ) ) === -1 || message.content.length === 1 || message.author.bot ) return;

  // gets the args of the command
  let args = message.content.slice( 1 ).split( ' ' );
  args = args.filter( n => n );

  // gets the command
  const commandName = args.shift().toLowerCase();

  const command = client.commands.get( commandName ) || client.commands.find( cmd => cmd.aliases && cmd.aliases.includes( commandName ) );

  if ( !command ) {
    return message.reply( 'Det finns inget sådant kommando.' );
  }

  if ( command.serverOnly && message.channel.type !== 'text' ) {
    return message.reply( 'Denna kommando finns bara tillgänglig i Digital Ungdom servern.' );
  }

  try {
    command.execute( message, args );
  } catch ( error ) {
    message.reply( `En error uppstog med commandot ${message.content}` );
    console.error( error );
  }
} );

// login to Discord with your app's token

MongoClient.connect( process.env.DB_URL, { useNewUrlParser: true }, async function ( err, mongoClient ) {
  if ( err ) return console.log( 'mongodb', err );
  global.db = mongoClient.db( 'digitalungdom' );
  client.login( process.env.BOT_TOKEN );
} );