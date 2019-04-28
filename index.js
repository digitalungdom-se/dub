/* global base_dir abs_path client memberProcess guild config version*/

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

global.memberProcess = {};

global.live = ( new Date() ).toISOString().slice( 0, 10 );
global.lastUpdated = '2019-04-28';

require( 'dotenv' ).config();
global.config = require( './config.json' );
global.version = ( require( './package.json' ) ).version;

const Discord = require( 'discord.js' );
const Spotify = require( 'node-spotify-api' );
const fs = require( 'fs-extra' );
const MongoClient = require( 'mongodb' ).MongoClient;
const Canvas = require( 'canvas' );
const fetch = require( 'node-fetch' );



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
  const channel = guild.channels.find( ch => ch.name === 'general' );
  channel.send( `startar botten på version: **${version}**.` );
} );

client.on( 'guildMemberAdd', async function ( member ) {
  const channel = member.guild.channels.find( ch => ch.name === 'general' );
  if ( !channel ) return;

  const canvas = Canvas.createCanvas( 800, 250 );
  const ctx = canvas.getContext( '2d' );

  const background = await Canvas.loadImage( abs_path( 'public/imgs/code.png' ) );
  ctx.drawImage( background, 0, 0, background.width, background.height );

  ctx.strokeStyle = '#ffffff';
  ctx.strokeRect( 0, 0, canvas.width, canvas.height );

  ctx.font = '25px sans-serif';
  ctx.fillStyle = '#000000';
  ctx.fillText( 'Välkommen till Digital Ungdom Servern,', 250, 75 );

  let font = 70;
  do {
    ctx.font = `${ font -= 5}px sans-serif`;
  } while ( ctx.measureText( `${member.displayName}!` ).width > canvas.width - 300 );

  ctx.fillStyle = '#000000';
  ctx.fillText( `${member.displayName}!`, 250, 140 );

  ctx.beginPath();
  ctx.arc( 125, 125, 100, 0, Math.PI * 2, true );
  ctx.closePath();
  ctx.clip();

  const buffer = await ( await fetch( member.user.displayAvatarURL ) ).buffer();

  const avatar = await Canvas.loadImage( buffer );
  ctx.drawImage( avatar, 25, 25, 200, 200 );

  const attachment = new Discord.Attachment( canvas.toBuffer(), 'välkommen.png' );

  return channel.send( `Välkommen till servern, ${member}!`, attachment );
} );

// mesage functions
client.on( 'message', async function ( message ) {
  if ( message.content === '$join' ) {
    return client.emit( 'guildMemberAdd', message.member || await message.guild.fetchMember( message.author ) );
  }

  if ( ( config.prefix.indexOf( message.content.charAt( 0 ) ) === -1 || message.content.length === 1 || message.author.bot ) && !memberProcess[ message.author.id ] ) return;
  else if ( memberProcess[ message.author.id ] ) {
    try {
      return client.commands.get( 'member' ).execute( message, message.content );
    } catch ( error ) {
      message.reply( `En error uppstog med commandot ${message.content}` );
      return console.error( error );
    }
  }

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

  if ( command.adminOnly ) {
    if ( !message.member.roles.find( r => r.name === 'admin' ) ) return message.reply( `du har inte behörighet att köra \`${command.name}\`.` );
  }

  if ( command.group === 'music' ) {
    if ( !message.member.voiceChannel ) return message.reply( 'Du måste vara i en ljud kanal.' );
    if ( guild.voiceConnection && ( message.member.voiceChannelID !== guild.me.voiceChannelID ) ) {
      message.delete();
      return message.reply( 'du och botten måste vara i samma ljudkanal' );
    }
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