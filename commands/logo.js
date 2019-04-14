module.exports = {
  name: 'logo',
  description: 'Skickar den nuvarande server logotyp',
  aliases: ['logotyp'],
  group: 'misc',
  usage: 'logo',
  execute(message, args) {
    const logoURL = 'https://cdn.discordapp.com/icons/468044152970674176/12781169588db1e946b3734cde72101c.webp?';

    message.reply('Här kommer serverns nuvarande logotyp',{file: logoURL});
  },
};
