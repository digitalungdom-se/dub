/* global  guild */

const axios = require( 'axios' );

module.exports = {
  name: 'status',
  description: 'Hämtar statusen av Digital Ungdom',
  aliases: [],
  group: 'digitalungdom',
  usage: 'status',
  example: 'status',
  serverOnly: false,
  adminOnly: false,
  async execute( message, args ) {
    const data = [];
    const status = ( await axios.get( 'https://digitalungdom.se/api/status' ) ).data;
    const embed = {
      'title': '__**STATUS**__',
      'color': 4086462,
      'timestamp': ( new Date() ).toISOString(),
      'fields': []
    };

    const boardField = { 'name': '__**Styrelse:**__', 'value': '' };
    const board = Object.keys( status.board );
    board.forEach( member => boardField.value += `**${member.charAt(0).toUpperCase() + member.slice(1)}**: ${status.board[member]}\n` );
    embed.fields.push( boardField );

    const serversField = { 'name': '__**Servers:**__', 'value': '' };
    const servers = Object.keys( status.server );
    servers.forEach( server => serversField.value += `**${server.charAt(0).toUpperCase() + server.slice(1)}**: ${status.server[server]}\n` );
    embed.fields.push( serversField );


    const membersField = { 'name': '__**Members:**__', 'value': '' };
    membersField.value = `**digitalungdom.se**: ${status.members['amount']}\n**discord**: ${guild.memberCount}`;
    embed.fields.push( membersField );


    message.reply( 'här kommer statusen av Digital Ungdom:', { 'embed': embed } );
  },
};