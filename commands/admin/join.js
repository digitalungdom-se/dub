/* global client guild */

module.exports = {
  name: 'join',
  description: 'Simulerar att en medlem joinar',
  longDescription: 'Simulerar att en medlem joinar',
  aliases: [],
  group: 'admin',
  usage: 'join <(@user)>',
  example: 'join @Ippyson#6200',
  serverOnly: true,
  adminOnly: true,
  async execute( message, args ) {
    let user = await message.mentions.members.first();
    if ( !user ) user = await guild.fetchMember( message.author );

    return client.emit( 'guildMemberAdd', user );
  },
};