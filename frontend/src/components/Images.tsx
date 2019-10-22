import React, { useEffect, useState } from "react";
import { RouteComponentProps, withRouter } from "react-router";
import {
    GridList, GridListTile,
    Theme, makeStyles, createStyles,
} from "@material-ui/core";

type Props = RouteComponentProps<{}>;

interface ImageResponse {
    image_url: string
}

const useStyles = makeStyles((theme: Theme) => {
    return createStyles({
        root: {
            display: "flex",
            flexWrap: "wrap",
            justifyContent: "space-around",
            overflow: "hidden",
            backgroundColor: theme.palette.background.paper,
        },
        gridList: {
            width: "100%",
        },
    });
});

const Images: React.FC<Props> = ({ history }) => {
    const classes = useStyles();
    const [images, setImages] = useState<ImageResponse[]>([]);
    useEffect(() => {
        (async () => {
            try {
                const res = await fetch("/api/index");
                if (!res.ok) {
                    if (res.status === 401) {
                        history.push("/");
                        return;
                    } else {
                        throw new Error(res.statusText);
                    }
                }
                const data: ImageResponse[] = await res.json();
                setImages(data);
            } catch (err) {

                window.console.error(err.message);
            }
        })();
    }, [history]);
    const tiles = images.map((image: ImageResponse, index: number) => {
        return (
          <GridListTile key={index}>
            <img src={image.image_url} alt="" />
          </GridListTile>
        );
    });
    return (
      <React.Fragment>
        <h2>Images</h2>
        <div className={classes.root}>
          <GridList cellHeight={120} className={classes.gridList} cols={10}>
            {tiles}
          </GridList>
        </div>
      </React.Fragment>
    );
};

export default withRouter(Images);
