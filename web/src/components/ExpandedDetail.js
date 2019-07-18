import React from 'react'

import HelmIcon from '../assets/helm_icon.png'
import KubeIcon from '../assets/kubernetes_icon.png'
import DataDogIcon from '../assets/data_dog_icon.png'
import TravisIcon from '../assets/travis_icon.png'
import SumoIcon from '../assets/sumo_logic_icon.png'
import Typography from '@material-ui/core/Typography'
import Switch from '@material-ui/core/Switch'
import MuiThemeProvider from '@material-ui/core/styles/MuiThemeProvider'
import createMuiTheme from '@material-ui/core/styles/createMuiTheme'
import FormControlLabel from '@material-ui/core/FormControlLabel'
import Paper from '@material-ui/core/Paper'
import ListItem from '@material-ui/core/ListItem'
import ListItemText from '@material-ui/core/ListItemText'
import List from '@material-ui/core/List'
import { CircularProgress, Link, Collapse } from '@material-ui/core'
import ExpandLess from '@material-ui/icons/ExpandLess'
import ExpandMore from '@material-ui/icons/ExpandMore'
import axios from 'axios'
import urljoin from 'url-join'

const theme = createMuiTheme({
  palette: {
    primary: {
      main: '#FF6612'
    },
    secondary: {
      main: '#FEAF0A'
    }
  }
})

const linkTheme = createMuiTheme({
  palette: {
    primary: {
      main: '#000000'
    },
    secondary: {
      main: '#FEAF0A'
    }
  }
})

const nestedStyle = {
  paddingLeft: theme.spacing(4)
}

const imgAlt = 'not found'

class ExpandedDetail extends React.Component {

  constructor(props) {
    super(props)
    this.state = {
      podsVisible: false,
      resourcesVisible: false,
      pods: null,
      resources: null
    }
  }

  podClick = () => {
    this.setState({podsVisible: !this.state.podsVisible})
    axios.get(urljoin(this.props.API_ADDRESS, 'releases', this.props.data.name, 'resources')).then(response => {
      this.setState({
        pods: response.data.filter((d) => d.kind === 'Pod'),
        resources: response.data.filter((d) => d.kind !== 'Pod')
      })
    })
  }

  resourceClick = () => {
    this.setState({resourcesVisible: !this.state.resourcesVisible})
  }

  render() {
    return (
      <div>
        {
          this.props.data ?
            <div>
              <div className="flex-container">
                <div className="helm-status">
                  <div className="right-padded">
                    <img src={HelmIcon} width="32" height="32" alt={imgAlt} />
                  </div>
                  <div>
                    <Typography variant="h6">{this.props.data.artifacts.chart.path}@{this.props.data.artifacts.chart.version}</Typography>
                  </div>
                </div>

                <MuiThemeProvider theme={theme}>
                  <div className="switch-status">
                    <FormControlLabel
                      control={
                        <Switch color="primary" checked={this.props.data.autoDeploy} />
                      }
                      label="AutoDeploy"
                    />
                  </div>
                </MuiThemeProvider>
                <MuiThemeProvider theme={linkTheme}>
                  <div className="dataDog-status">
                    <div className="right-padded">
                      <img src={DataDogIcon} width="32" height="32" alt={imgAlt} />
                    </div>
                    <div>
                      <Typography variant="h6">
                        <Link href={this.props.data.monitoring.datadog.dashboard}>Dashboard</Link>
                      </Typography>
                    </div>
                    <div>
                      <Typography variant="h6">&nbsp;|&nbsp;</Typography>
                    </div>
                    <div>
                      <Typography variant="h6">
                        <Link href={this.props.data.monitoring.datadog.monitors}>Monitor</Link>
                      </Typography>
                    </div>
                  </div>
                </MuiThemeProvider>
              </div>
              <div className="flex-container">
                <div className="kube-status">
                  <div className="right-padded">
                    <img src={KubeIcon} width="32" height="32" alt={imgAlt} />
                  </div>
                  <div>
                    <Typography variant="h6">Resources</Typography>
                  </div>
                </div>
                <div className="log-status">
                  <div className="right-padded">
                    <img src={SumoIcon} width="32" height="32" alt={imgAlt} />
                  </div>
                  <div>
                    <Typography variant="h6">
                      <Link href={this.props.data.monitoring.sumologic}>Logs</Link>
                    </Typography>
                  </div>
                  <div className="double-padded">
                    <img src={TravisIcon} width="32" height="32" alt={imgAlt} />
                  </div>
                  <div>
                    <Typography variant="h6">
                      <Link href={this.props.data.build.travis}>Build</Link>
                    </Typography>
                  </div>
                </div>
              </div>
              <div className="kube-resource-list">
                <Paper style={{ maxHeight: '200px', minWidth: '100%', overflow: 'auto' }}>
                  <List
                    component="nav"
                  >
                    <ListItem button onClick={this.podClick}>
                      <ListItemText>Pods</ListItemText>
                      {this.state.podsVisible ? <ExpandLess /> : <ExpandMore />}
                    </ListItem>
                    <Collapse in={this.state.podsVisible} timeout='auto' unmountOnExit>
                      <List component="div" disablePadding>
                        {
                          this.state.pods ?
                          <div>
                            {
                              this.state.pods.map((pod) =>
                                <ListItem button style={nestedStyle} key={pod.metadata.name}>
                                  <ListItemText>{pod.metadata.name}</ListItemText>
                                </ListItem>
                              )
                            }
                          </div>
                          : <CircularProgress/>
                        }
                      </List>
                    </Collapse>
                    <ListItem button onClick={this.resourceClick}>
                      <ListItemText>Other Resources</ListItemText>
                      {this.state.resourcesVisible ? <ExpandLess /> : <ExpandMore />}
                    </ListItem>
                    <Collapse in={this.state.resourcesVisible} timeout='auto' unmountOnExit>
                      <List component="div" disablePadding>
                        {
                          this.state.resources ?
                          <div>
                            {
                              this.state.resources.map((resource) =>
                                <ListItem style={nestedStyle} key={resource.kind}>
                                  <ListItemText>{resource.kind}</ListItemText>
                                </ListItem>
                              )
                            }
                          </div> : <CircularProgress/>
                        }
                      </List>
                    </Collapse>
                  </List>
                </Paper>
              </div>
            </div> : <CircularProgress />
        }
      </div>
    )
  }
}

export default ExpandedDetail
