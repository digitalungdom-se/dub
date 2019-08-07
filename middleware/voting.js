/* global include voteDic */

const createVoteEmbed = include( 'utils/embeds/createVoteEmbed' );

module.exports.add = async function ( message, user ) {
  if ( user.bot ) return;
  if ( voteDic[ message.message.id ] ) {
    const id = message.message.id;
    const index = message.emoji.name;
    const userID = user.id;

    const allowed = Array.from( message.message.reactions.keys() );

    if ( allowed.length !== voteDic[ id ].options.length ) return message.emoji.reaction.remove( user.id );
    else {
      voteDic[ id ][ 'users' ][ userID ] = index;

      const embed = createVoteEmbed( voteDic[ id ].embed, voteDic[ id ].options, voteDic[ id ].message )[ 0 ];

      await voteDic[ id ].message.edit( { 'embed': embed } );
    }
  }
};

module.exports.remove = async function ( message, user ) {
  if ( user.bot ) return;
  const id = message.message.id;
  const index = message.emoji.name;
  const userID = user.id;

  if ( voteDic[ message.message.id ] ) {
    delete voteDic[ id ][ 'users' ][ userID ];

    const [ embed, reactions ] = createVoteEmbed( voteDic[ id ].embed, voteDic[ id ].options, voteDic[ id ].message );
    await voteDic[ id ].message.edit( { 'embed': embed } );
  }
};