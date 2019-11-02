import React, { useEffect, useState } from "react";
import { useHistory, useLocation } from "react-router";
import { Link as RouterLink, LinkProps } from "react-router-dom";
import {
    Card, CardActionArea, CardMedia, CardContent,
    Box, CircularProgress, Grid, Link, Typography,
    Theme, makeStyles, createStyles,
} from "@material-ui/core";

import SearchBox from "./SearchBox";
import { ImageResponse } from "../common/interfaces";

const useStyles = makeStyles((theme: Theme) => {
    return createStyles({
        progress: {
            margin: theme.spacing(2),
        },
        card: {
            maxWidth: 148,
        },
        media: {
            height: 148,
            width: 148,
        },
    });
});

const Images: React.FC = () => {
    const classes = useStyles();
    const history = useHistory();
    const location = useLocation();
    const [images, setImages] = useState<ImageResponse[]>([]);
    const [loading, setLoading] = useState<boolean>(false);
    useEffect(() => {
        if (location.search.length === 0) {
            return;
        }
        const params = new URLSearchParams(location.search);
        params.set("limit", "100");
        setImages([]);
        setLoading(true);
        fetch(`/api/images?${params}`).then((res: Response) => {
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
        }).finally(() => {
            setLoading(false);
        });
    }, [history, location]);
    const cards = images.map((image: ImageResponse) => {
        const link = React.forwardRef<HTMLAnchorElement, Omit<LinkProps, "to">>(
            (props, ref) => <RouterLink innerRef={ref} to={`/image/${image.id}`} {...props} />,
        );
        return (
          <Box key={image.id} m={0.25}>
            <Card className={classes.card}>
              <Link component={link}>
                <CardActionArea>
                  <CardMedia
                      className={classes.media}
                      image={image.image_url} />
                </CardActionArea>
              </Link>
              <CardContent>
                <Typography variant="body2">
                  {image.size}×{image.size}
                </Typography>
                <Typography variant="body2" color="textSecondary">
                  {image.label_name}
                </Typography>
              </CardContent>
            </Card>
          </Box>
        );
    });
    const progress = (
      <Grid container justify="center">
        <CircularProgress className={classes.progress} />
      </Grid>
    );
    return (
      <React.Fragment>
        <h2>Images</h2>
        <SearchBox />
        {loading && progress}
        <Box display="flex" flexWrap="wrap" mt={2}>
          {cards}
        </Box>
      </React.Fragment>
    );
};

export default Images;
