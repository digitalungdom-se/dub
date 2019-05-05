/* global searchMessage searchList guild include controller */

const ytdl = require( 'ytdl-core' );

const createMusicEmbed = include( 'utils/embeds/createMusicEmbed' );

const reactions = { '0⃣': 0, '1⃣': 1, '2⃣': 2, '3⃣': 3, '4⃣': 4 };

module.exports = async function ( message, user ) {
  if ( user.bot ) return;
  if ( searchMessage.id !== message.message.id ) return;
  message.emoji.reaction.remove( user.id );

  const member = await guild.fetchMember( user );
  if ( !member.voiceChannel || guild.voiceConnection && ( member.voiceChannelID !== guild.me.voiceChannelID ) ) {
    return;
  }

  const url = `https://www.youtube.com${searchList[ reactions[message.emoji.name] ].url}`;

  const metadata = await ytdl.getBasicInfo( url );
  metadata.user = message.message.author.id;
  const embed = createMusicEmbed( metadata );

  controller.add( { 'url': url, 'embed': embed, 'metadata': metadata }, member );

  global.searchList = false;
  if ( searchMessage ) {
    searchMessage.delete();
    global.searchMessage = false;
  }
};