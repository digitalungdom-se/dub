/* global include */

const version = (include('package.json')).version;

module.exports = {
  name: 'version',
  description: 'Anger versionen av botten',
  aliases: [],
  group: 'misc',
  usage: 'version',
  execute(message, args) {
    message.reply(`Den nuvarande versionen av botten är: **${version}**`);
  },
};
