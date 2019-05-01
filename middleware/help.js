/* global help include */

const helpEmbeds = include( 'utils/embeds/helpEmbeds' );

module.exports = async function ( message, user ) {
  if ( user.bot ) return;
  if ( !help[ user.id ] || help[ user.id ].id !== message.message.id ) return;

  const helpMessage = help[ user.id ];

  const page = message.emoji.name;

  if ( page === helpEmbeds.reactions[ helpEmbeds.reactions.length - 1 ] ) return helpMessage.delete();

  const embed = helpEmbeds[ page ]();

  return helpMessage.edit( { 'embed': embed } );
};