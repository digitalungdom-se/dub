/* global base_dir abs_path */

global.base_dir = __dirname;
global.abs_path = function ( path ) {
  return `${base_dir}/${path}`;
};
global.include = function ( file ) {
  return require( abs_path( file ) );
};

require( 'dotenv' ).config();
const config = require( './config.json' );

const Discord = require( 'discord.js' );
const fs = require( 'fs-extra' );



// create a new Discord client
const client = new Discord.Client();
client.commands = new Discord.Collection();

const commandFiles = fs.readdirSync( './commands' ).filter( file => file.endsWith( '.js' ) );

for ( const file of commandFiles ) {
  const command = require( `./commands/${file}` );
  client.commands.set( command.name, command );
}

// when the client is ready, run this code
// this event will only trigger one time after logging in
client.once( 'ready', () => {
  console.log( 'Ready!' );
} );

client.on( 'message', message => {
  if ( message.channel.type !== 'text' ) {
    return message.reply( 'Jag finns bara tillgänglig i digital ungdom servern.' );
  }

  if ( config.prefix.indexOf( message.content.charAt( 0 ) ) === -1 || message.content.length === 1 || message.author.bot ) return;

  // gets the args of the command
  const args = message.content.slice( 1 ).split( ' ' );
  // gets the command
  const commandName = args.shift().toLowerCase();

  const command = client.commands.get( commandName ) || client.commands.find( cmd => cmd.aliases && cmd.aliases.includes( commandName ) );

  if ( !command ) {
    return message.reply( 'Det finns inget sådant kommando.' );
  }

  try {
    command.execute( message, args );
  } catch ( error ) {
    message.reply( `En error uppstog med commandot ${message.content}` );
    console.error( error );
  }
} );

// login to Discord with your app's token
client.login( process.env.BOT_TOKEN );