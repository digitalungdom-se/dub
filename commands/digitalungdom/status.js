const axios = require( 'axios' );

module.exports = {
  name: 'status',
  description: 'Hämtar statusen av Digital Ungdom',
  aliases: [],
  group: 'digitalungdom',
  usage: 'status',
  serverOnly: false,
  async execute( message, args ) {
    const data = [];
    const status = ( await axios.get( 'https://digitalungdom.se/api/status' ) ).data;
    data.push( 'här kommer statusen av Digital Ungdom:' );
    data.push( '__**Styrelse:**__' );
    const board = Object.keys( status.board );
    board.forEach( member => data.push( `**${member.charAt(0).toUpperCase() + member.slice(1)}**: ${status.board[member]}` ) );

    data.push( '\n__**Servers:**__' );
    const servers = Object.keys( status.server );
    servers.forEach( server => data.push( `**${server}**: ${status.server[server]}` ) );

    data.push( '\n__**Members:**__' );
    const members = Object.keys( status.members );
    members.forEach( member => data.push( `**${member}**: ${status.members[member]}` ) );

    message.reply( data, { split: true } );
  },
};