/* global base_dir abs_path client guild version lastUpdated live */
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
  global.lastUpdated = '2019-05-01';

  global.config = require( './config.json' );
  global.version = ( require( './package.json' ) ).version;

  client.user.setActivity( 'Kelvin\'s cat', { type: 'WATCHING' } );
  guild.channels.find( ch => ch.name === 'general' ).send( `startar boten på version: **${version}**.`, {
    'embed': {
      'description': '__**INFORMATION OM BOTEN**__',
      'color': 4086462,
      'fields': [ {
          'name': 'VERSION',
          'value': `${version} (${lastUpdated})`,
          'inline': true
        },
        {
          'name': 'LIVE SEDAN',
          'value': `${live}`,
          'inline': true
        },
        {
          'name': 'KÄLLKOD',
          'value': '[github](https://github.com/kelszo/dub)',
          'inline': true
        },
        {
          'name': 'MEDARBETARE',
          'value': '<@217632464531619852>, <@228889878861971456>',
          'inline': true
        }
      ]
    }
  } );

  console.log( 'Ready!' );
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

// message handler middleware
client.on( 'message', messageHandler );

MongoClient.connect( process.env.DB_URL, { useNewUrlParser: true }, async function ( err, mongoClient ) {
  if ( err ) return console.log( 'mongodb', err );
  global.db = mongoClient.db( 'digitalungdom' );
  client.login( process.env.BOT_TOKEN );
} );