/* global controller client guild include */

// '⏮'

const Controller = include( 'utils/controller' );

module.exports = {
  name: 'controller',
  description: 'Skickar musik kontrollen.',
  longDescription: 'Skickar en kontroll vilket man kan kontrollera musiken med.',
  aliases: [ 'kontroll' ],
  group: 'music',
  usage: 'controller',
  example: 'controller',
  serverOnly: true,
  adminOnly: false,
  async execute( message, args ) {
    if ( controller.queue === 0 || !controller || !guild.me.voiceChannelID ) global.controller = new Controller( client, guild );
    else controller.newController();
  },
};