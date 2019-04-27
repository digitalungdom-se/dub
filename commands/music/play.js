/* global musicQueue player musicVolume client searchList searchMessage include */

const ytdl = require( 'ytdl-core' );
const ytSearch = require( 'yt-search' );
const util = require( 'util' );

const ytSearchAsync = util.promisify( ytSearch );
const createMusicEmbed = include( 'utils/createMusicEmbed' );

async function play( connection, message ) {
  const url = musicQueue[ 0 ].url;
  global.player = connection.playStream( ytdl( url, { quality: 'highestaudio', filter: 'audioonly' } ), { bitrate: 192000 } );
  player.setVolume( musicVolume );

  const metadata = await ytdl.getBasicInfo( url );
  let seconds = metadata.length_seconds % 60;
  if ( seconds < 10 ) seconds = `0${seconds}`;

  const embed = createMusicEmbed( metadata );
  embed.description = `Spelas nu på begäran av <@${musicQueue[ 0 ].queuer.id}>. Den är ${Math.floor(metadata.length_seconds/60)}:${seconds} minuter lång.`;
  embed.content = 'Spelar nu:';

  global.musicQueue[ 0 ].started = new Date();

  client.user.setActivity( metadata.title );

  message.channel.send( { embed } );

  player.on( 'end', function () {
    global.musicQueue.shift();
    if ( musicQueue.length !== 0 ) {
      play( connection, message );
    } else {
      connection.disconnect();
      client.user.setActivity( 'Kelvin\'s cat', { type: 'WATCHING' } );
    }
  } );
}

module.exports = {
  name: 'play',
  description: 'Spelar länken i din nuvarande kanal',
  aliases: [ 'spela', 'pl' ],
  group: 'music',
  usage: 'play <youtube link|search term|@user>',
  example: 'play https://www.youtube.com/watch?v=dQw4w9WgXcQ',
  serverOnly: true,
  adminOnly: false,
  async execute( message, args ) {
    message.delete();
    let url = '';
    if ( args.length === 0 ) return message.reply( 'Du måste välja en låt.' );
    if ( musicQueue.length > 1000 ) return message.reply( 'Kön är full' );
    if ( message.mentions.members.first() ) {
      const user = ( message.mentions.members.first() );
      const game = user.presence.game;
      if ( game && game.name === 'Spotify' ) {
        // state is the artist and details is the song name
        const song = `${game.details} ${game.state}`;
        const result = await ytSearchAsync( song );
        url = `https://www.youtube.com${result.videos[ 0 ].url}`;
      } else {
        return message.reply( `<@${user.user.id}> spelar ingen låt på spotify.` );
      }
    } else if ( ytdl.validateURL( args[ 0 ] ) ) {
      url = args[ 0 ];
    } else if ( searchList && ( args[ 0 ] >= 0 || args[ 0 ] <= 5 ) ) {
      const index = Math.floor( args[ 0 ] );
      url = `https://www.youtube.com${searchList[ index ].url}`;
    } else {
      global.searchList = false;
      if ( searchMessage ) {
        client.channels.get( searchMessage.channelId ).fetchMessage( searchMessage.id ).then( msg => msg.delete() );
        global.searchMessage = false;
      }
      const data = [];
      data.push( 'här är låtarna du kan välja mellan:' );
      global.searchList = ( await ytSearchAsync( args.join( ' ' ) ) ).videos.slice( 0, 5 );

      searchList.forEach( function ( song, index ) {
        data.push( `**${index}.** ${song.author.name} | ${song.title}` );
      } );
      data.push( '\n Välj låt med `$play <n>` där n är låt siffran.' )
      const msg = await message.reply( data, { split: true } );
      global.searchMessage = { id: msg.id, channelId: msg.channel.id };
      return;
    }
    const metadata = await ytdl.getBasicInfo( url );

    const embed = createMusicEmbed( metadata );

    musicQueue.push( { url: url, title: `${metadata.author.name} | ${metadata.title}`, embed: embed, lengthSeconds: metadata.length_seconds, queuer: message.author } );
    if ( !message.guild.voiceConnection ) {
      const connection = await message.member.voiceChannel.join();
      play( connection, message );
    } else {
      let seconds = metadata.length_seconds % 60;
      if ( seconds < 10 ) seconds = `0${seconds}`;
      const minutes = Math.floor( metadata.length_seconds / 60 );
      embed.description = `Lades till i kön på begäran av ${message.author}. Den är ${minutes}:${seconds} minuter lång.`;

      message.channel.send( { embed } );
    }
    global.searchList = false;
    if ( searchMessage ) {
      client.channels.get( searchMessage.channelId ).fetchMessage( searchMessage.id ).then( msg => msg.delete() );
      global.searchMessage = false;
    }
  },
};