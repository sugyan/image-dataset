import React, { useRef, useEffect, useState } from "react";
import { RouteComponentProps } from "react-router";
import { Paper } from "@material-ui/core";

import { ImageResponse } from "../common/interfaces";

type Props = RouteComponentProps<{ id: string }>;

const ImageViewer: React.FC<Props> = ({ match, history }) => {
    const [images, setImages] = useState<ImageResponse[]>([]);
    const canvas = useRef<HTMLCanvasElement>(null);
    if (canvas && canvas.current) {
        const ctx = canvas.current.getContext("2d")!;
        const image = new Image();
        image.onload = () => {
            ctx.drawImage(image, 0, 0, 512, 512);
        };
        image.src = images[0].image_url;
    }
    useEffect(() => {
        const params: URLSearchParams = new URLSearchParams({ key: match.params.id });
        fetch(
            `/api/images?${params.toString()}`
        ).then((res: Response) => {
            if (res.ok) {
                return res.json();
            }
            if (res.status === 401) {
                history.push("/");
                return;
            }
            throw new Error(res.statusText);
        }).then((images: ImageResponse[]) => {
            setImages(images);
        }).catch((err: Error) => {
            window.console.error(err.message);
        });
    }, [history, match.params.id]);
    return (
      <React.Fragment>
        <h2>Image</h2>
        <Paper>
          <canvas height={512} width={512} ref={canvas} />
        </Paper>
      </React.Fragment>
    );
};

export default ImageViewer;
