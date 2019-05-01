/* global client searchList searchMessage include controller */

const ytdl = require( 'ytdl-core' );
const ytSearch = require( 'yt-search' );
const util = require( 'util' );

const ytSearchAsync = util.promisify( ytSearch );
const createMusicEmbed = include( 'utils/embeds/createMusicEmbed' );

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
    let url = '';
    if ( args.length === 0 ) return message.reply( 'Du måste välja en låt.' ).then( msg => { msg.delete( 10000 ); } );
    if ( controller.queue.length > 30 ) return message.reply( 'Kön är full' ).then( msg => { msg.delete( 10000 ); } );
    if ( message.mentions.members.first() ) {
      const user = ( message.mentions.members.first() );
      const game = user.presence.game;
      if ( game && game.name === 'Spotify' ) {
        // state is the artist and details is the song name
        const song = `${game.details} ${game.state}`;
        const result = await ytSearchAsync( song );
        url = `https://www.youtube.com${result.videos[ 0 ].url}`;
      } else {
        return message.reply( `<@${user.user.id}> spelar ingen låt på spotify.` ).then( msg => { msg.delete( 10000 ); } );
      }
    } else if ( ytdl.validateURL( args[ 0 ] ) ) {
      url = args[ 0 ];
    } else if ( searchList && ( args[ 0 ] >= 0 || args[ 0 ] <= 5 ) ) {
      const index = Math.floor( args[ 0 ] );
      url = `https://www.youtube.com${searchList[ index ].url}`;
    } else {
      global.searchList = false;
      if ( searchMessage ) {
        searchMessage.delete();
        global.searchMessage = false;
      }
      const embed = {
        'title': '__**Här är låtarna du kan välja mellan:**__',
        'description': '',
        'color': 4086462
      };
      global.searchList = ( await ytSearchAsync( args.join( ' ' ) ) ).videos.slice( 0, 5 );

      searchList.forEach( function ( song, index ) {
        embed.description += `**${index}.** ${song.author.name} | ${song.title}\n`;
      } );
      embed.description += '\n Välj låt med `$play <n>` där n är låt siffran.';

      const msg = await message.reply( { 'embed': embed } );
      global.searchMessage = msg;
      return;
    }
    const metadata = await ytdl.getBasicInfo( url );
    metadata.user = message.author.id;
    const embed = createMusicEmbed( metadata );

    controller.add( { 'url': url, 'embed': embed, 'metadata': metadata }, message );

    global.searchList = false;
    if ( searchMessage ) {
      searchMessage.delete();
      global.searchMessage = false;
    }
  },
};