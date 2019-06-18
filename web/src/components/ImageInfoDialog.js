import React from 'react'
import Dialog from '@material-ui/core/Dialog'
import { DialogContent, DialogActions, Button, DialogTitle, Typography } from '@material-ui/core';

class ImageInfoDialog extends React.Component {

    getRegistry() {
        var arr = this.props.docker.image.split('/')
        return arr[0]
    }

    getRepo() {
        var arr = this.props.docker.image.split('/')
        return arr[1]
    }

    getTag() {
        return this.props.docker.tag
    }

    getURI() {
        return this.props.docker.image + ':' + this.getTag()
    }

    render() {
        return (
            <Dialog open={this.props.open} onClose={this.props.handleClose} maxWidth={true}>
                <DialogTitle>Docker Image Information</DialogTitle>
                <DialogContent>
                    <Typography>Docker Registry: {this.getRegistry()}</Typography>
                    <Typography>Repository: {this.getRepo()}</Typography>
                    <Typography>Tag: {this.getTag()}</Typography>
                    <Typography>Full URI: {this.getURI()}</Typography>
                </DialogContent>
                <DialogActions>
                    <Button onClick={this.props.handleClose}>Close</Button>
                </DialogActions>
            </Dialog>
        )
    }
}

export default ImageInfoDialog;