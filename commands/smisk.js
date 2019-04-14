module.exports = {
  name: 'slap',
  description: 'Smiskar en styrelsemedlem.',
  aliases: ['smisk'],
  group: 'misc',
  usage: 'slap <medlem>',
  execute(message, args) {
    const name = args[0].toLowerCase();

    const members = {
      'kelvin': '<@217632464531619852>',
      'simon': '<@228889878861971456>',
      'douglas': '<@!297671552823066626>',
    };

    message.channel.send(`Du har varit en riktig stygg pojk ${members[name]}.`)
  },
};
