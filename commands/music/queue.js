/* global musicQueue */

module.exports = {
  name: 'queue',
  description: 'Kollar vilka låtar som är i kön.',
  aliases: [ 'kö', 'q' ],
  group: 'music',
  usage: 'queue',
  serverOnly: true,
  execute( message, args ) {
    if ( musicQueue.length === 0 ) return message.reply( 'det finns inga låtar i kön.' );
    const data = [];
    data.push( 'dessa låtar är i kön:' );
    musicQueue.forEach( function ( song, index ) {
      data.push( `**${index}.** ${song.title}` );
    } );
    message.reply( data, { split: true } );
  },
};