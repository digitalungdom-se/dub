/* global controller guild */

const reactions = [ '❎', '⏯', '⏭', '➕', '➖' ];

module.exports = async function ( message, user ) {
  if ( user.bot ) return;
  if ( !controller || controller.message.id !== message.message.id ) return;
  message.emoji.reaction.remove( user.id );

  const member = await guild.fetchMember( user );
  if ( !member.voiceChannel || guild.voiceConnection && ( member.voiceChannelID !== guild.me.voiceChannelID ) ) {
    return;
  }
  const command = message.emoji.name;

  switch ( command ) {
  case '❎':
    controller.stop();
    break;
  case '⏯':
    controller.pauseResume();
    break;
  case '⏭':
    controller.skip();
    break;
  case '➕':
    controller.setVolume( { inc: 0.1 } );
    break;
  case '➖':
    controller.setVolume( { inc: -0.1 } );
    break;
  default:
    return;
  }
};