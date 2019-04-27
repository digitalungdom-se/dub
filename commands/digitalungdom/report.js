/* global db client include */

const getUserByDiscordId = include( 'models/get' ).getUserByDiscordId;

module.exports = {
  name: 'report',
  description: 'Anmäl någon. Skicka via DM för att utföra delvis anonym.',
  aliases: [ 'anmäl', 're' ],
  group: 'digitalungdom',
  usage: 'report <@user> <reason>',
  example: 'report @Zigolox#0919 Han är taskig :(',
  serverOnly: false,
  adminOnly: false,
  async execute( message, args ) {
    if ( message.channel.type === 'text' ) message.delete();
    let reportedUser;

    if ( message.channel.type === 'text' ) {
      if ( !message.mentions.members.first() ) return message.reply( 'du måste @ vem du vill reporta.' );
      if ( args.length < 2 ) return message.reply( 'du måste ge en kort anledning.' );

      reportedUser = message.mentions.members.first().user;
    } else if ( message.channel.type === 'dm' ) {
      if ( args.length < 4 ) return message.reply( 'du måste ange vem du vill reporta och ge en kort anledning.' );
      let reported = args[ 0 ];
      if ( reported.startsWith( '@' ) ) reported = reported.slice( 1 );
      if ( /#\d{4}$/.test( reported ) ) reported = reported.slice( 0, -5 );

      reportedUser = client.users.find( user => user.username == reported );
      if ( !reportedUser ) return message.reply( 'Kunde inte hitta användare, se till så att det inte är deras nickname och att alla versaler är korrekta.' );
    }

    let reason = args;
    reason.shift();
    reason = reason.join( ' ' );

    const reportedId = reportedUser.id;
    const reportedUsername = reportedUser.username;

    const authorId = message.author.id;
    const authorUsername = message.author.username;

    let reportedDUId, authorDUId;
    [ reportedDUId, authorDUId ] = await Promise.all( [
      getUserByDiscordId( reportedId ),
      getUserByDiscordId( authorId ),
    ] );

    if ( reportedDUId ) reportedDUId = reportedDUId._id;
    if ( authorDUId ) authorDUId = authorDUId._id;

    await db.collection( 'notifications' ).insertOne( {
      'type': 'report',
      'where': 'discord',
      'message': reason,
      'reported': {
        'id': reportedDUId,
        'discordId': reportedId,
        'discordUsername': reportedUsername
      },
      'author': {
        'id': authorDUId,
        'discordId': authorId,
        'discordUsername': authorUsername
      }
    } );

    return message.author.send( 'Tack för din medverkan, vi kommer inom en snar framtid att granska din anmälning.' );
  },
};