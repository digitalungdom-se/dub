/* global musicQueue spotify include */

const createMusicEmbed = include( 'utils/createMusicEmbed' );
const msToTimeStamp = include( 'utils/msToTimeStamp' );

module.exports = {
  name: 'playing',
  description: 'Visar upp vilken låt som spelas nu.',
  aliases: [ 'np', 'spelas' ],
  group: 'music',
  usage: 'playing',
  serverOnly: true,
  async execute( message, args ) {
    if ( message.mentions.members.first() ) {
      const user = ( message.mentions.members.first() );
      const game = user.presence.game;
      if ( game && game.name === 'Spotify' ) {
        const metadata = {};
        metadata.author = {};
        // state is the artist and details is the song name
        metadata.title = game.details;
        metadata.author.name = game.state;

        // gets current time stamp
        const timeStamp = msToTimeStamp( new Date() - game.timestamps.start );

        const resultSpotify = await spotify.search( { type: 'track', query: `${game.details} ${game.state}`, limit: 1 } );
        const authorId = resultSpotify.tracks.items[ 0 ].album.artists[ 0 ].id;

        // gets thumbnail
        metadata.thumbnail_url = resultSpotify.tracks.items[ 0 ].album.images[ 0 ].url;
        metadata.author.avatar = ( await spotify.request( `https://api.spotify.com/v1/artists/${authorId}` ) ).images[ 0 ].url;

        const embed = createMusicEmbed( metadata );
        embed.description = timeStamp;

        return message.reply( `<@${user.user.id}> lyssnar på:`, { embed } );
      } else return message.reply( `<@${user.user.id}> spelar ingen låt på spotify.` );
    } else if ( message.guild.voiceConnection ) {
      if ( musicQueue.length === 0 ) message.reply( 'det spelas inget.' );
      else {
        const data = musicQueue[ 0 ];
        const embed = data.embed;
        const currentTimeStamp = msToTimeStamp( new Date() - data.started );
        const totalTimeStap = msToTimeStamp( data.lengthSeconds * 1000 );
        embed.description = `Spelas nu på begäran av ${data.queuer}. Det har gått ${currentTimeStamp}/${totalTimeStap}`;
        message.reply( { embed } );
      }
    } else message.reply( 'det spelas inget.' );
  },
};