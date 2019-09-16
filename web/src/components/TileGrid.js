/* eslint react/prop-types: 0 */
/* eslint react/jsx-no-bind: 0 */

import React from 'react'
import PropTypes from 'prop-types'
import SingleGridCell from './GridCell'

class ReactExpandableGrid extends React.Component {
  
  constructor(props) {
    super(props)

    this.state = {
      expanded: false,
      selected_id: '',
      gridData: this.props.gridData
    }
  }

  findIdNode(target) { // takes event target and returns the id (written to make sure the expanded detail renders from the correct reference point as in the original component each tile on the grid was just an img tag)
    let t = target
    while (t != null) {
      if (String(t.className) === "SingleGridCell") {
        return parseInt(t.id.substring(10))
      }
      t = t.parentNode
    }
    return -1
  }
  
  renderExpandedDetail(target) {
    let thisIdNumber = this.findIdNode(target)
    let detail = document.getElementById('expandedDetail')

    let ol = document.getElementById("grid_cell_" + thisIdNumber.toString()).parentNode
    let lengthOfList = parseInt(ol.childNodes.length)
    let startingIndex = thisIdNumber + 1

    let insertedFlag = false

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
      let detail = document.getElementById('expandedDetail')
      detail.style.display = 'none'
    })
  }

  handleCellClick(event) {
    let target = event.target
    if (this.state.expanded) {
      if (this.state.selected_id === event.target.id) { // Clicking on already opened detail
        this.closeExpandedDetail()
        this.renderExpandedDetail(target)
      } else { // Clicking on a different thumbnail, when detail is already expanded
        this.setState({
          expanded: true,
          selected_id: event.target.id
        }, function afterStateChange() {
          let detail = document.getElementById('expandedDetail')

          this.renderExpandedDetail(target)

          detail.style.display = 'block'
        })
      }
    } else {
      this.setState({
        expanded: true,
        selected_id: event.target.id
      }, function afterStateChange() {
        let detail = document.getElementById('expandedDetail')

        this.renderExpandedDetail(target)

        detail.style.display = 'block'
      })
    }
  }

  generateGrid() {
    let grid = []
    let idCounter = -1 // To help simplify mapping to object array indices. For example, <li> with 0th id corresponds to 0th child of <ol>
    let gridData = this.state.gridData
    for (var i in gridData) {
      idCounter = idCounter + 1
      let thisUniqueKey = 'grid_cell_' + idCounter.toString()
      if (this.matchesQuery(gridData[i], this.props.query)) {
        grid.push(<SingleGridCell API_ADDRESS={this.props.API_ADDRESS} handleCellClick={this.handleCellClick.bind(this)} key={thisUniqueKey} id={thisUniqueKey} cellMargin={this.props.cellMargin} SingleGridCellData={gridData[i]} cellSize={this.props.cellSize} />)
      }
    }

    return grid
  }

  matchesQuery(release, query) {
    const fields = [release.name, release.owner.squad, release.status];
    return fields
      .map(f => f.toLowerCase())
      .some(f => f.includes(query));
  }

  render() {
    let rows = this.generateGrid()

    const cssForGridDetailExpansion = {
      width: '100%',
      position: 'relative'
    }

    const cssForGridList = {
      listStyle: 'none',
      padding: 0,
      display: 'inline-block',
      width: '100%'
    }

    const cssForTheGridHolder = {
      width: '100%',
      backgroundColor: this.props.bgColor,
      margin: 0,
      textAlign: 'center'
    }
    return (
      <div id='GridDetailExpansion' style={cssForGridDetailExpansion}>
        <div id='theGridHolder' style={cssForTheGridHolder}>
          <ol id='gridList' style={cssForGridList}>
            {rows}
          </ol>
        </div>
      </div>
    )
  }
}

ReactExpandableGrid.propTypes = {
  gridData: PropTypes.array,
  cellSize: PropTypes.number,
  cellMargin: PropTypes.number,
  bgColor: PropTypes.string,
}

ReactExpandableGrid.defaultProps = {
  gridData: [],
  cellSize: 250,
  cellMargin: 25,
  bgColor: '#f2f2f2',
}

export default ReactExpandableGrid
