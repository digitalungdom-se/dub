/* global status client controller */

module.exports = {
  name: 'dubstatus',
  description: 'Ändrar statusen av boten.',
  aliases: [],
  group: 'admin',
  usage: 'dubstatus <PLAYING|STREAMING|LISTENING|WATCHING> <status>',
  example: 'dubstatus WATCHING youtube',
  serverOnly: true,
  adminOnly: true,
  execute( message, args ) {
    if ( [ 'PLAYING', 'STREAMING', 'LISTENING', 'WATCHING' ].indexOf( args[ 0 ] ) === -1 ) return message.reply( `okänt aktivitets typ: \`${args[0]}\`` );
    const type = args[ 0 ];

    let acitivity = args;
    acitivity.shift();
    if ( acitivity.length === 0 ) return message.reply( 'du måste ge en status.' );
    acitivity = acitivity.join( ' ' );

    global.status = { 'acitivity': acitivity, 'type': type };

    if ( !controller.playing ) client.user.setActivity( status.acitivity, { 'type': status.type } );

    return message.reply( `sätter nu statusen till \`${type} ${acitivity}\`` ).then( msg => { msg.delete( 10000 ); } );
  },
};