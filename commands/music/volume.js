/* global player musicVolume */

module.exports = {
  name: 'volume',
  description: 'Ändrar volumen av botten eller visar den nuvarande volym.',
  aliases: [ 'volym', 'vol' ],
  group: 'music',
  usage: 'volume <0-2 || [current , nuvarande]>',
  example: 'volume 0.1',
  serverOnly: true,
  adminOnly: false,
  execute( message, args ) {
    if ( args.length === 0 ) return message.reply( `den nuvarande volymen är ${player.volume}` );
    if ( !player ) return message.reply( 'botten måste vara aktiv för att ändra volym.' );
    if ( [ 'current', 'nuvarande' ].indexOf( args[ 0 ] ) !== -1 ) return message.reply( `den nuvarande volymen är ${player.volume}` );
    if ( args[ 0 ] < 0 || args[ 0 ] > 2 ) return message.reply( 'volymnivån måste ligga mellan 0 och 2.' );
    global.musicVolume = args[ 0 ];
    player.setVolume( musicVolume );
  },
};