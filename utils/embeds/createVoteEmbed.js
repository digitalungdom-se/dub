/* global voteDic */

const emojis = [ '0⃣', '1⃣', '2⃣', '3⃣', '4⃣', '5⃣', '6⃣', '7⃣', '8⃣', '9⃣', '🔟' ];

module.exports = function createVoteEmbed( embed, options, message ) {
  let reactions = [];
  embed.fields = [];

  const score = {};
  if ( voteDic[ message.id ] ) {
    for ( const user of Object.keys( voteDic[ message.id ][ 'users' ] ) ) {
      if ( score[ voteDic[ message.id ][ 'users' ][ user ] ] === undefined ) score[ voteDic[ message.id ][ 'users' ][ user ] ] = 1;
      else score[ voteDic[ message.id ][ 'users' ][ user ] ] += 1;
    }
  }
  for ( const [ index, option ] of options.entries() ) {
    if ( option.length > 1024 ) return message.reply( 'ett val får max vara 1024 tecken.' );
    let votingNumber = emojis[ index ];
    reactions.push( votingNumber );

    let amount;
    let percentage;

    if ( voteDic[ message.id ] && Object.keys( voteDic[ message.id ][ 'users' ] ).length > 0 && score[ votingNumber ] ) {
      amount = score[ votingNumber ];
      percentage = Math.trunc( ( score[ votingNumber ] / Object.keys( voteDic[ message.id ][ 'users' ] ).length ) * 100 );
    } else {
      amount = 0;
      percentage = 0;
    }

    embed.addField( `${votingNumber}: ${option}`, `${amount} (${percentage}%)` );
  }

  return [ embed, reactions ];
};