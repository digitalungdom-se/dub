/* global musicVolume controller */

module.exports = {
  name: 'volume',
  description: 'Ändrar musikens volym.',
  aliases: [ 'volym', 'vol' ],
  group: 'music',
  usage: 'volume <0-2 | [current , nuvarande]>',
  example: 'volume 0.1',
  serverOnly: true,
  adminOnly: false,
  execute( message, args ) {
    if ( args.length === 0 ) return message.reply( `den nuvarande volymen är ${controller.volume}` ).then( msg => { msg.delete( 10000 ); } );
    if ( !controller.player ) return message.reply( 'boten måste vara aktiv för att ändra volym.' ).then( msg => { msg.delete( 10000 ); } );
    if ( [ 'current', 'nuvarande' ].indexOf( args[ 0 ] ) !== -1 ) return message.reply( `den nuvarande volymen är ${controller.volume}` ).then( msg => { msg.delete( 10000 ); } );
    if ( args[ 0 ] < 0 || args[ 0 ] > 2 ) return message.reply( 'volymnivån måste ligga mellan 0 och 2.' ).then( msg => { msg.delete( 10000 ); } );

    controller.setVolume( { set: args[ 0 ] } );
  },
};