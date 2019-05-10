/* global base_dir abs_path client guild */
require( 'dotenv' ).config();

global.base_dir = __dirname;
global.abs_path = function ( path ) {
  return `${base_dir}/${path}`;
};
global.include = function ( file ) {
  return require( abs_path( file ) );
};

const Discord = require( 'discord.js' );
const Spotify = require( 'node-spotify-api' );
const fs = require( 'fs-extra' );
const MongoClient = require( 'mongodb' ).MongoClient;

const Controller = require( './utils/controller' );

const votingMiddleware = require( './middleware/voting' );
const welcomerMiddleware = require( './middleware/welcomer' );
const helpMiddleware = require( './middleware/help' );
const controllerMiddleware = require( './middleware/controller' );
const messageHandler = require( './middleware/messageHandler' );
const searchListMiddleware = require( './middleware/searchList' );

// create a new Discord client
global.client = new Discord.Client();
client.commands = new Discord.Collection();

// create spotify API client
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
client.once( 'ready', async function () {
  global.guild = client.guilds.get( process.env.GUILD_ID );

  const musicChannel = guild.channels.find( ch => ch.name === 'music' );
  let messages = await musicChannel.fetchMessages( { limit: 100 } );
  while ( messages.size ) {
    musicChannel.bulkDelete( messages );
    messages = await musicChannel.fetchMessages( { limit: 100 } );
  }

  guild.channels.find( ch => ch.name === 'voting' ).send( 'Alla omröstningar innan detta meddelande är nu stängda.' );
  global.controller = new Controller( client, guild );
  global.searchList = false;
  global.searchMessage = false;

  global.memberProcess = {};

  global.voteDic = {};

  global.help = {};

  global.cooldown = {
    verify: {}
  };

  global.live = ( new Date() ).toISOString().slice( 0, 10 );
  global.lastUpdated = '2019-05-10';

  global.config = require( './config.json' );
  global.version = ( require( './package.json' ) ).version;

  global.status = { 'acitivity': 'Kelvin\'s cat', 'type': 'WATCHING' };

  client.user.setActivity( status.acitivity, { 'type': status.type } );

  console.log( 'Ready!', ( new Date() ).toISOString() );
} );

// welcomer middleware
client.on( 'guildMemberAdd', welcomerMiddleware );

// voting middleware
client.on( 'messageReactionAdd', votingMiddleware.add );
client.on( 'messageReactionRemove', votingMiddleware.remove );

// help middleware
client.on( 'messageReactionAdd', helpMiddleware );
client.on( 'messageReactionRemove', helpMiddleware );

// music controller middleware
client.on( 'messageReactionAdd', controllerMiddleware );

// search list middleware
client.on( 'messageReactionAdd', searchListMiddleware );

// message handler middleware
client.on( 'message', messageHandler );

// error handler
process.on( 'unhandledRejection', error => console.error( 'Uncaught Promise Rejection', ( new Date() ).toISOString(), error ) );
process.on( 'uncaughtException', error => console.error( 'Uncaught Exception', ( new Date() ).toISOString(), error ) );
client.on( 'error', error => console.error( 'Uncaught Error', ( new Date() ).toISOString(), error ) );

MongoClient.connect( process.env.DB_URL, { useNewUrlParser: true }, async function ( err, mongoClient ) {
  if ( err ) return console.log( 'mongodb', err );
  global.db = mongoClient.db( 'digitalungdom' );
  client.login( process.env.BOT_TOKEN );
} );