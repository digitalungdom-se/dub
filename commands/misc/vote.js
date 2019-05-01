/* global voteDic include guild */

const RichEmbed = require( 'discord.js' ).RichEmbed;
const createVoteEmbed = include( 'utils/embeds/createVoteEmbed' );

module.exports = {
  name: 'vote',
  description: 'Börja en röstning om en fråga',
  aliases: [ 'votestart', 'rösta', 'börjarösta' ],
  group: 'misc',
  usage: 'votestart "<titel>" "<alternativ 1>" "<alternativ 2>" "<alternativ 3>"',
  example: 'votestart "vad ska vi käka?" "sushi" "taccos" "pizza"',
  serverOnly: true,
  async execute( message, args ) {
    if ( !message.deleted ) message.delete( 10000 );
    let options = args.join( ' ' ).split( '"' ).filter( ( e, i ) => i % 2 === 1 );
    const title = options[ 0 ];
    options.shift();
    if ( options.length > 11 ) return message.reply( 'du får max ha 11 val.' );
    if ( options.length < 2 ) return message.reply( 'du måste ha några val.' );

    let embed = new RichEmbed()
      .setTitle( title )
      .setAuthor( message.author.username, message.author.displayAvatarURL )
      .setColor( 4086462 )
      .setTimestamp();

    let reactions;
    [ embed, reactions ] = createVoteEmbed( embed, options, message );

    message.reply( `grattis du har på börjat en röstningen "*${title}*"! Gå till \`voting\` kanalen för att se den` );

    const channel = guild.channels.find( ch => ch.name === 'voting' );
    if ( !channel ) return;

    const reactionMessage = await channel.send( { 'embed': embed } );

    const id = reactionMessage.id;
    voteDic[ id ] = { options };
    voteDic[ id ][ 'users' ] = {};
    voteDic[ id ][ 'score' ] = {};
    voteDic[ id ][ 'message' ] = reactionMessage;
    voteDic[ id ][ 'embed' ] = embed;

    for ( const reaction of reactions ) {
      await reactionMessage.react( reaction );
    }

    return;
  },
};