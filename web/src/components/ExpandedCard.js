import React from 'react'
import { Dialog, DialogContent, DialogActions, List, Button, DialogTitle, ListItem, ListItemText } from '@material-ui/core';
import ExpandLess from '@material-ui/icons/ExpandLess'
import ExpandMore from '@material-ui/icons/ExpandMore'
import { Collapse } from '@material-ui/core'
import { MuiThemeProvider } from "@material-ui/core/styles"
import Typography from '@material-ui/core/Typography'
import Link from '@material-ui/core/Link'
import FormControlLabel from '@material-ui/core/FormControlLabel'
import Switch from '@material-ui/core/Switch'
import axios from 'axios'
import urljoin from 'url-join'

import * as themes from '../Themes'
import * as constants from '../Constants'

import HelmIcon from '../assets/helm_icon.png'
import DataDogIcon from '../assets/data_dog_icon.png'

class ExpandedCard extends React.Component {
    constructor(props) {
        super(props)
	    this.state = {
	    	podsVisible: false,
		resourceString: ""
	    }
    }

    podClick = () => {
        this.setState({ podsVisible: !this.state.podsVisible })
        axios.get(urljoin(this.props.API_ADDRESS, 'releases', this.props.data.name, 'resources')).then(response => {
            this.setState({ resourceString: response.data.status })
        })
    }

    render() {
        return (
            <Dialog open={this.props.open} onClose={this.props.handleClose} maxWidth="xl">
                <DialogTitle>{this.props.data.name}</DialogTitle>
                <DialogContent>
                    <div className="flex-container" style={{width:700}}>
                        <div className="helm-status">
                            <div className="right-padded">
                                <img src={HelmIcon} width="32" height="32" alt={constants.imgAlt} />
                            </div>
                            <div className="center">
                                <Typography variant="h7">{this.props.data.artifacts.chart.path}@{this.props.data.artifacts.chart.version}</Typography>
                            </div>
                        </div>

                        <MuiThemeProvider theme={themes.standard}>
                            <div className="switch-status">
                                <FormControlLabel
                                    control={
                                        <Switch color="primary" checked={this.props.data.autoDeploy} />
                                    }
                                    label="AutoDeploy"
                                />
                            </div>
                        </MuiThemeProvider>

                        <MuiThemeProvider theme={themes.link}>
                            <div className="dataDog-status">
                                <div className="right-padded">
                                    <img src={DataDogIcon} width="32" height="32" alt={constants.imgAlt}/>
                                </div>
                                <div className="center">
                                    <Typography variant="h7">
                                        <Link href={this.props.data.monitoring.datadog.dashboard}>Dashboard</Link>
                                    </Typography>
                                </div>
                                <div className="center">
                                    <Typography variant="h7">&nbsp;|&nbsp;</Typography>
                                </div>
                                <div className="center">
                                    <Typography variant="h7">
                                        <Link href={this.props.data.monitoring.datadog.monitors}>Monitor</Link>
                                    </Typography>
                                </div>
                            </div>
                        </MuiThemeProvider>
                    </div>
                    <List component="nav">
                        <ListItem button onClick={this.podClick}>
                            <ListItemText>Helm Status</ListItemText>
                            {this.state.podsVisible ? <ExpandLess /> : <ExpandMore />}
                        </ListItem>
                        <Collapse in={this.state.podsVisible} timeout='auto' unmountOnExit>
                            <pre>{this.state.resourceString}</pre>
                        </Collapse>
                    </List>
                </DialogContent>
                <DialogActions>
                    <Button onClick={this.props.handleClose}>Close</Button>
                </DialogActions>
            </Dialog>
        )
    }
}

export default ExpandedCard
