/* eslint react/prop-types: 0 */
/* eslint react/jsx-no-bind: 0 */

import React from 'react'
import PropTypes from 'prop-types'
import Card from '@material-ui/core/Card'
import CardContent from '@material-ui/core/CardContent'
import Typography from '@material-ui/core/Typography'
import Chip from '@material-ui/core/Chip'
import DoneIcon from '@material-ui/icons/Done'
import ExpandIcon from '@material-ui/icons/ExpandMore'
import CompressIcon from '@material-ui/icons/ExpandLess'
import IconButton from '@material-ui/core/IconButton'
import Green from '@material-ui/core/colors/green';
import SelectionDialog from './SelectionDialog'
import CircularProgress from '@material-ui/core/CircularProgress'

import TimePassed from '../assets/time_passed.png'
import SlackIcon from '../assets/slack_icon.png'
import GitIcon from '../assets/octocat.png'
import DockerIcon from '../assets/docker_icon.png'
import KubeIcon from '../assets/kubernetes_icon.png'
import JenkinsIcon from '../assets/jenkins_icon.png'
import TravisIcon from '../assets/travis_icon.png'
import ExpandedDetail from './ExpandedDetail';

import { createMuiTheme, MuiThemeProvider } from '@material-ui/core/styles';

const deployTagTheme = createMuiTheme({
  palette: {
    primary: {
      main: '#4caf50'
    },
    secondary: {
      main: '#f44336'
    },
    default: {
      main: '#9e9e9e'
    }
  }
})

var keys = []

class SingleGridCell extends React.Component {
  constructor(props) {
    super(props)

    this.state = {
      expanded: false,
      selected_id: '',
      window_width: window.innerWidth,
      repoSelector: false,
    }
  }

  cellClick(event) {
    this.props.handleCellClick(event)
  }

  render() {
    var SingleGridCellStyle = {
      backgroundSize: this.props.cellSize,
      width: this.props.cellSize,
      height: this.props.cellSize,
      display: 'inline-block',
      margin: this.props.cellMargin,
      marginBottom: 25,
      position: 'relative'
    }
    console.log(keys)
    var cardStyle = {
      width: this.props.cellSize,
      height: this.props.cellSize
    }
    return (
      <div style={SingleGridCellStyle} id={this.props.id} className='SingleGridCell'>
        <div>
          <Card style={cardStyle}>
            <CardContent >
              <Typography variant="h5" component="h2">
                ServiceName
              </Typography>
              <MuiThemeProvider theme={deployTagTheme}>
                <Chip
                  icon={<DoneIcon />}
                  label="Deployed"
                  color="primary"
                  variant="outlined"
                  clickable
                />
              </MuiThemeProvider>
              <div>
                <IconButton>
                  <img src={TimePassed} />
                </IconButton>
                5 min ago
              </div>
              <div className='row-align'>
                <SelectionDialog />
                <IconButton>
                  <img src={SlackIcon} width="32" height="32" />
                </IconButton>
                <IconButton>
                  <img src={DockerIcon} width="32" height="32" />
                </IconButton>
              </div>
              <div>
                {
                  this.state.expanded // Not working
                    ?
                    <IconButton onClick={this.cellClick.bind(this)}>
                      <CompressIcon />
                    </IconButton>
                    :
                    <IconButton onClick={this.cellClick.bind(this)}>
                      <ExpandIcon />
                    </IconButton>
                }
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    )
  }
}

class ReactExpandableGrid extends React.Component {

  constructor(props) {
    super(props)

    this.state = {
      expanded: false,
      selected_id: '',
      gridData: JSON.parse(this.props.gridData)
    }
  }

  findIdNode(target) { // takes event target and returns the id
    var t = target
    var found = false
    while (!found && t != null) {
      console.log(t.className)
      if (String(t.className) === "SingleGridCell") {
        found = true;
        return parseInt(t.id.substring(10))
      }
      t = t.parentNode
    }
    return -1;
  }

  renderExpandedDetail(target) {
    //var thisId = target.id
    //var className = target.className;
    var thisIdNumber = this.findIdNode(target)
    var detail = document.getElementById('expandedDetail')

    console.log(detail)
    var ol = document.getElementById("grid_cell_" + thisIdNumber.toString()).parentNode
    var lengthOfList = parseInt(ol.childNodes.length)
    var startingIndex = thisIdNumber + 1

    var insertedFlag = false

    ol.insertBefore(detail, ol.childNodes[lengthOfList])

    for (var i = startingIndex; i < lengthOfList; i++) {
      if (ol.childNodes[i].className === 'SingleGridCell') {
        if (ol.childNodes[i].offsetTop !== ol.childNodes[thisIdNumber].offsetTop) {
          ol.childNodes[i].insertAdjacentElement('beforebegin', detail)
          insertedFlag = true
          break
        }
      }
    }

    if (insertedFlag === false) {
      ol.childNodes[lengthOfList - 1].insertAdjacentElement('afterend', detail)
    }
  }

  closeExpandedDetail() {
    this.setState({
      expanded: false,
      selected_id: ''
    }, function afterStateChange() {
      var detail = document.getElementById('expandedDetail')
      detail.style.display = 'none'
    })
  }

  handleCellClick(event) {
    var target = event.target
    var thisIdNumber = this.findIdNode(target)
    var className = event.target.className
    console.log(event)
    console.log(thisIdNumber);
    if (this.state.expanded) { // expanded == true
      if (this.state.selected_id === event.target.id) { // Clicking on already opened detail
        this.closeExpandedDetail()
        this.renderExpandedDetail(target)
      } else { // Clicking on a different thumbnail, when detail is already expanded
        this.setState({
          expanded: true,
          selected_id: event.target.id
        }, function afterStateChange() {
          var detail = document.getElementById('expandedDetail')

          this.renderExpandedDetail(target)

          detail.style.display = 'block'
        })
      }
    } else { // expanded == false
      this.setState({
        expanded: true,
        selected_id: event.target.id
      }, function afterStateChange() {
        var detail = document.getElementById('expandedDetail')

        this.renderExpandedDetail(target)

        detail.style.display = 'block'
      })
    }
  }

  generateGrid() {
    var grid = []
    var idCounter = -1 // To help simplify mapping to object array indices. For example, <li> with 0th id corresponds to 0th child of <ol>
    var gridData = this.state.gridData
    var k = []
    for (var i in gridData) {
      idCounter = idCounter + 1
      var thisUniqueKey = 'grid_cell_' + idCounter.toString()
      k.push(thisUniqueKey)
      grid.push(<SingleGridCell handleCellClick={this.handleCellClick.bind(this)} key={thisUniqueKey} id={thisUniqueKey} cellMargin={this.props.cellMargin} SingleGridCellData={gridData[i]} cellSize={this.props.cellSize} />)
    }
    keys = k
    var cssforExpandedDetail = {
      backgroundColor: this.props.detailBackgroundColor,
      height: this.props.detailHeight,
      display: 'none',
      position: 'relative',
      padding: '20px',
      transition: 'display 2s ease-in-out 0.5s'
    }

    var cssforExpandedDetailImage = {
      display: 'inline-block',
      maxWidth: this.props.ExpandedDetail_image_size,
      width: '100%',
      height: 'auto',
      align: 'center',
      position: 'absolute',
      top: 0,
      bottom: 0,
      left: 0,
      right: 0,
      margin: 'auto'
    }

    var cssforExpandedDetailTitle = {
      backgroundColor: this.props.ExpandedDetail_title_bgColor,
      width: '100%',
      height: 'auto',
      marginBottom: '15px'
    }

    var cssforExpandedDetailDescription = {
      backgroundColor: this.props.ExpandedDetail_description_bgColor,
      color: this.props.ExpandedDetail_font_color,
      width: 'auto%',
      height: '80%',
      marginRight: '30px',
      marginLeft: '30px',
      textAlign: 'justify'
    }

    var cssforExpandedDetailLeft
    var cssforExpandedDetailRight

    cssforExpandedDetailLeft = {
      width: this.props.ExpandedDetail_left_width,
      height: '100%',
      float: 'left',
      position: 'relative'
    }

    cssforExpandedDetailRight = {
      width: this.props.ExpandedDetail_right_width,
      height: '100%',
      float: 'right',
      position: 'relative'
    }

    var cssForDescriptionLink = {
      textDecoration: 'none',
      position: 'relative',
      float: 'bottom',
      bottom: 20,
      cursor: 'pointer'
    }

    var cssForImageLink = {
      cursor: 'pointer'
    }

    var cssforExpandedDetailClose = {
      textDecoration: 'none',
      position: 'relative',
      float: 'right',
      top: 10,
      right: 10,
      cursor: 'pointer'
    }

    // Make Mobile Friendly
    if (window.innerWidth < this.props.show_mobile_style_from_width) {
      cssforExpandedDetailLeft = {
        width: '0%',
        height: '100%',
        float: 'left',
        position: 'relative',
        display: 'none'
      }

      cssforExpandedDetailRight = {
        width: '100%',
        height: '100%',
        float: 'right',
        position: 'relative'
      }
    }

    var closeX
    if (this.props.ExpandedDetail_closeX_bool) {
      closeX = 'X'
    } else {
      closeX = ''
    }

    grid.push( // Expanded Detail here
      <li style={cssforExpandedDetail} key='expandedDetail' id='expandedDetail'>
        <ExpandedDetail />
      </li>
    )

    return grid
  }

  render() {
    var rows = this.generateGrid()

    var cssForGridDetailExpansion = {
      width: '100%',
      position: 'relative'
    }

    var cssForGridList = {
      listStyle: 'none',
      padding: 0,
      display: 'inline-block'
    }

    var cssForTheGridHolder = {
      width: '100%',
      backgroundColor: this.props.bgColor,
      margin: 0,
      textAlign: 'center'
    }

    // var cssForSelectedArrow = {
    //   width: 0,
    //   height: 0,
    //   borderLeft: '20px solid transparent',
    //   borderRight: '20px solid transparent',
    //   borderBottom: '30px solid' + this.props.detailBackgroundColor,
    //   marginTop: this.props.cellSize,
    //   marginLeft: this.props.cellSize / 2 - 20,
    //   display: 'none'
    // }

    return (
      <div id='GridDetailExpansion' style={cssForGridDetailExpansion}>
        <div id='theGridHolder' style={cssForTheGridHolder}>
          <ol id='gridList' style={cssForGridList}>
            {rows}
          </ol>
        </div>
        {/*<div id='selected_arrow' style={cssForSelectedArrow} />*/}
      </div>
    )
  }
}

ReactExpandableGrid.propTypes = {
  gridData: PropTypes.string,
  cellSize: PropTypes.number,
  cellMargin: PropTypes.number,
  bgColor: PropTypes.string,
  detailWidth: PropTypes.string, // in %
  detailHeight: PropTypes.number,
  detailBackgroundColor: PropTypes.string,
  ExpandedDetail_right_width: PropTypes.string, // in %
  ExpandedDetail_left_width: PropTypes.string, // in %
  ExpandedDetail_description_bgColor: PropTypes.string,
  ExpandedDetail_title_bgColor: PropTypes.string,
  ExpandedDetail_img_bgColor: PropTypes.string,
  ExpandedDetail_link_text: PropTypes.string,
  ExpandedDetail_font_color: PropTypes.string,
  ExpandedDetail_closeX_bool: PropTypes.bool,
  show_mobile_style_from_width: PropTypes.number
}

var data = [
  { 'img': 'http://i.imgur.com/zIEjP6Q.jpg', 'link': 'https://www.instagram.com/p/BRFjVZtgSJD/', 'title': 'Westland Tai Poutini National Park', 'description': 'Photo by @christopheviseux / The Westland Tai Poutini National Park in New Zealand’s South Island offers a remarkable opportunity to take a guided walk on a glacier. A helicopter drop high on the Franz Josef Glacier, provides access to explore stunning ice formations and blue ice caves. Follow me for more images around the world @christopheviseux #newzealand #mountain #ice' },
  { 'img': 'http://i.imgur.com/rCrvQTv.jpg', 'link': 'https://www.instagram.com/p/BQ6_Wa2gmdR/', 'title': 'Dubai Desert Conservation Reserve', 'description': 'Photo by @christopheviseux / Early morning flight on a hot air balloon ride above the Dubai Desert Conservation Reserve. Merely an hour drive from the city, the park was created to protect indigenous species and biodiversity. The Arabian Oryx, which was close to extinction, now has a population well over 100. There are many options to explore the desert and flying above may be one of the most mesmerizing ways. Follow me @christopheviseux for more images from the Middle East. #dubai #desert' },
  { 'img': 'http://i.imgur.com/U8iVzVl.jpg', 'link': 'https://www.instagram.com/p/BQyfDiKAEq9/', 'title': 'Crumbling Reflections', 'description': 'Photo @pedromcbride // Crumbling Reflections: Much has changed in Cuba over the 17 years I have visited this island. But much has stayed the same. Time still ticks at a Cuban pace and old cars still run… I don’t know how... and while pockets of new construction and renovation exist thanks to a growing tourism boom, most buildings are crumbling and cracking under the Caribbean climate. But amidst the hardship, nostalgia and messy vitality, the Cuban people keep moving, like their cars. And somehow, they do it with a colorful friendliness and warmth that always amazes me. To see more, follow @pedromcbride #cuba #havana #photo #workshop @natgeoexpeditions #reflection #photooftheday #petemcbride.' },
  { 'img': 'http://i.imgur.com/Ky9aJlE.jpg', 'link': 'https://www.instagram.com/p/BQxf6CEgD8p/', 'title': 'Impalas', 'description': 'Impetious young impala go head-to-head as they practice sparring. A talent they will need later in life when the rut begins. Photographed on assignment for @natgeotravel in Kruger National Park. For more images from Kruger, South Africa, follow @kengeiger #natgeotravel #krugernationalpark' },
  { 'img': 'http://i.imgur.com/mf3qfzt.jpg', 'link': 'https://www.instagram.com/p/BQvy7gbgynF/', 'title': 'Elephants', 'description': 'Photo by @ronan_donovan // Two bull African elephants at dawn in Uganda\'s Murchison Falls National Park. See more from Uganda with @ronan_donovan.' },
  { 'img': 'http://i.imgur.com/zIEjP6Q.jpg', 'link': 'https://www.instagram.com/p/BRFjVZtgSJD/', 'title': 'Westland Tai Poutini National Park', 'description': 'Photo by @christopheviseux / The Westland Tai Poutini National Park in New Zealand’s South Island offers a remarkable opportunity to take a guided walk on a glacier. A helicopter drop high on the Franz Josef Glacier, provides access to explore stunning ice formations and blue ice caves. Follow me for more images around the world @christopheviseux #newzealand #mountain #ice' },
  { 'img': 'http://i.imgur.com/rCrvQTv.jpg', 'link': 'https://www.instagram.com/p/BQ6_Wa2gmdR/', 'title': 'Dubai Desert Conservation Reserve', 'description': 'Photo by @christopheviseux / Early morning flight on a hot air balloon ride above the Dubai Desert Conservation Reserve. Merely an hour drive from the city, the park was created to protect indigenous species and biodiversity. The Arabian Oryx, which was close to extinction, now has a population well over 100. There are many options to explore the desert and flying above may be one of the most mesmerizing ways. Follow me @christopheviseux for more images from the Middle East. #dubai #desert' },
  { 'img': 'http://i.imgur.com/U8iVzVl.jpg', 'link': 'https://www.instagram.com/p/BQyfDiKAEq9/', 'title': 'Crumbling Reflections', 'description': 'Photo @pedromcbride // Crumbling Reflections: Much has changed in Cuba over the 17 years I have visited this island. But much has stayed the same. Time still ticks at a Cuban pace and old cars still run… I don’t know how... and while pockets of new construction and renovation exist thanks to a growing tourism boom, most buildings are crumbling and cracking under the Caribbean climate. But amidst the hardship, nostalgia and messy vitality, the Cuban people keep moving, like their cars. And somehow, they do it with a colorful friendliness and warmth that always amazes me. To see more, follow @pedromcbride #cuba #havana #photo #workshop @natgeoexpeditions #reflection #photooftheday #petemcbride.' },
  { 'img': 'http://i.imgur.com/Ky9aJlE.jpg', 'link': 'https://www.instagram.com/p/BQxf6CEgD8p/', 'title': 'Impalas', 'description': 'Impetious young impala go head-to-head as they practice sparring. A talent they will need later in life when the rut begins. Photographed on assignment for @natgeotravel in Kruger National Park. For more images from Kruger, South Africa, follow @kengeiger #natgeotravel #krugernationalpark' },
  { 'img': 'http://i.imgur.com/mf3qfzt.jpg', 'link': 'https://www.instagram.com/p/BQvy7gbgynF/', 'title': 'Elephants', 'description': 'Photo by @ronan_donovan // Two bull African elephants at dawn in Uganda\'s Murchison Falls National Park. See more from Uganda with @ronan_donovan.' },
  { 'img': 'http://i.imgur.com/zIEjP6Q.jpg', 'link': 'https://www.instagram.com/p/BRFjVZtgSJD/', 'title': 'Westland Tai Poutini National Park', 'description': 'Photo by @christopheviseux / The Westland Tai Poutini National Park in New Zealand’s South Island offers a remarkable opportunity to take a guided walk on a glacier. A helicopter drop high on the Franz Josef Glacier, provides access to explore stunning ice formations and blue ice caves. Follow me for more images around the world @christopheviseux #newzealand #mountain #ice' },
  { 'img': 'http://i.imgur.com/rCrvQTv.jpg', 'link': 'https://www.instagram.com/p/BQ6_Wa2gmdR/', 'title': 'Dubai Desert Conservation Reserve', 'description': 'Photo by @christopheviseux / Early morning flight on a hot air balloon ride above the Dubai Desert Conservation Reserve. Merely an hour drive from the city, the park was created to protect indigenous species and biodiversity. The Arabian Oryx, which was close to extinction, now has a population well over 100. There are many options to explore the desert and flying above may be one of the most mesmerizing ways. Follow me @christopheviseux for more images from the Middle East. #dubai #desert' },
  { 'img': 'http://i.imgur.com/rCrvQTv.jpg', 'link': 'https://www.instagram.com/p/BQ6_Wa2gmdR/', 'title': 'Dubai Desert Conservation Reserve', 'description': 'Photo by @christopheviseux / Early morning flight on a hot air balloon ride above the Dubai Desert Conservation Reserve. Merely an hour drive from the city, the park was created to protect indigenous species and biodiversity. The Arabian Oryx, which was close to extinction, now has a population well over 100. There are many options to explore the desert and flying above may be one of the most mesmerizing ways. Follow me @christopheviseux for more images from the Middle East. #dubai #desert' }
]

ReactExpandableGrid.defaultProps = {
  gridData: JSON.stringify(data),
  cellSize: 250,
  cellMargin: 25,
  bgColor: '#f2f2f2',
  detailWidth: '100%',
  detailHeight: 300,
  detailBackgroundColor: '#D9D9D9',
  ExpandedDetail_right_width: '60%',
  ExpandedDetail_left_width: '40%',
  ExpandedDetail_image_size: 300,
  ExpandedDetail_description_bgColor: '#D9D9D9',
  ExpandedDetail_title_bgColor: '#D9D9D9',
  ExpandedDetail_img_bgColor: '#D9D9D9',
  ExpandedDetail_link_text: '→ Link',
  ExpandedDetail_font_color: '#434343',
  ExpandedDetail_closeX_bool: true,
  show_mobile_style_from_width: 600,
}

export default ReactExpandableGrid
