/* global musicQueue */

module.exports = {
  name: 'queue',
  description: 'Kollar vilka låtar som är i kön.',
  aliases: [ 'kö', 'qu' ],
  group: 'music',
  usage: 'queue',
  serverOnly: true,
  adminOnly: false,
  execute( message, args ) {
    if ( musicQueue.length === 0 ) return message.reply( 'det finns inga låtar i kön.' );
    const data = [];
    data.push( '\n**Denna låt spelas nu:**' );
    data.push( `**${0}.** ${musicQueue[0].title}\n` );

    if ( musicQueue.length > 1 ) {
      data.push( '**Dessa låtar är i kön:**' );
      musicQueue.slice( 1 ).forEach( function ( song, index ) {
        data.push( `**${index + 1}.** ${song.title}` );
      } );
    }
    message.reply( data, { split: true } );
  },
};