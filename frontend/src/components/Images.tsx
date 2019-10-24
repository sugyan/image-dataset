import React, { useEffect, useState } from "react";
import { RouteComponentProps, withRouter } from "react-router";
import { Link } from "react-router-dom";
import {
    Card, CardActionArea, CardMedia, CardContent,
    Box, Typography,
    Theme, makeStyles, createStyles,
} from "@material-ui/core";

import { ImageResponse } from "../common/interfaces";

type Props = RouteComponentProps<{}>;

const useStyles = makeStyles((theme: Theme) => {
    return createStyles({
        card: {
            maxWidth: 148,
        },
        media: {
            height: 148,
            width: 148,
        },
    });
});

const Images: React.FC<Props> = ({ history }) => {
    const classes = useStyles();
    const [images, setImages] = useState<ImageResponse[]>([]);
    useEffect(() => {
        fetch(
            "/api/images",
        ).then((res: Response) => {
            if (res.ok) {
                return res.json();
            }
            if (res.status === 401) {
                history.push("/");
                return;
            }
            throw new Error(res.statusText);
        }).then((data: ImageResponse[]) => {
            setImages(data);
        }).catch((err: Error) => {
            window.console.error(err.message);
        });
    }, [history]);
    const cards = images.map((image: ImageResponse) => {
        return (
          <Box key={image.id} m={0.25}>
            <Card className={classes.card}>
              <CardActionArea>
                <Link to={`/image/${image.id}`}>
                  <CardMedia
                      className={classes.media}
                      image={image.image_url} />
                </Link>
                <CardContent>
                  <Typography variant="body2">
                    {image.size}×{image.size}
                  </Typography>
                  <Typography variant="body2" color="textSecondary">
                    {image.label_name}
                  </Typography>
                </CardContent>
              </CardActionArea>
            </Card>
          </Box>
        );
    });
    return (
      <React.Fragment>
        <h2>Images</h2>
        <Box display="flex" flexWrap="wrap">
          {cards}
        </Box>
      </React.Fragment>
    );
};

export default withRouter(Images);
