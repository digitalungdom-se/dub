/* global guild include db */

const getUserByDiscordId = include( 'models/get' ).getUserByDiscordId;
const createNotificationEmbed = include( 'utils/embeds/createNotificationEmbed' );

module.exports = {
  name: 'kick',
  description: 'Sparkar användaren från servern.',
  aliases: [ 'sparka' ],
  group: 'admin',
  usage: 'kick <@user> <reason>',
  example: 'kick @Ippyson#6200 han var taskig',
  serverOnly: true,
  adminOnly: true,
  async execute( message, args ) {
    if ( !message.mentions.members.first() ) return message.reply( 'du måste @ vem du vill kicka.' ).then( msg => { msg.delete( 10000 ); } );
    if ( args.length < 2 ) return message.reply( 'du måste ge en kort anledning.' ).then( msg => { msg.delete( 10000 ); } );

    let reason = args;
    reason.shift();
    reason = reason.join( ' ' );

    const member = await guild.fetchMember( message.mentions.members.first() );
    member.kick( reason );

    const kickedID = member.user.id;
    const kickedUsername = member.user.username;

    const adminID = message.author.id;
    const adminUsername = message.author.username;

    let kickedDUID, adminDUID;
    [ kickedDUID, adminDUID ] = await Promise.all( [
      getUserByDiscordId( kickedID ),
      getUserByDiscordId( adminID ),
    ] );

    if ( kickedDUID ) kickedDUID = kickedDUID._id;
    if ( adminDUID ) adminDUID = adminDUID._id;

    await db.collection( 'notifications' ).insertOne( {
      'type': 'kick',
      'where': 'discord',
      'message': reason,
      'kicked': {
        'id': kickedDUID,
        'discordID': kickedID,
        'discordUsername': kickedUsername
      },
      'admin': {
        'id': adminDUID,
        'discordID': adminID,
        'discordUsername': adminUsername
      }
    } );

    const notification = createNotificationEmbed( 'KICK', `sparkade <@${kickedID}>.\n\n **Anledning:** ${reason}`, 16711680, { 'id': adminID, 'name': kickedUsername, 'url': member.user.displayAvatarURL } );
    const notificationChannel = guild.channels.find( ch => ch.name === 'notifications' );
    notificationChannel.send( '@here, ny notifikation', { 'embed': notification } );

    return message.reply( { 'embed': notification } );
  },
};