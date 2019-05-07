/* global help include guild */

const helpEmbeds = include( 'utils/embeds/createHelpEmbeds' );

module.exports = async function ( message, user ) {
  if ( user.bot ) return;
  if ( !help[ user.id ] || help[ user.id ].id !== message.message.id ) return;

  const helpMessage = help[ user.id ];

  const page = message.emoji.name;

  if ( page === '🔥' ) return helpMessage.delete();
  if ( page === '🚨' && !( await guild.fetchMember( user ) ).roles.find( r => r.name === 'admin' ) ) return;

  const embed = helpEmbeds[ page ]();

  return helpMessage.edit( { 'embed': embed } );
};