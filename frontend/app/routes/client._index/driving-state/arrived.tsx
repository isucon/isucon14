import { Form, useNavigate } from "@remix-run/react";
import { MouseEventHandler, useCallback, useRef, useState } from "react";
import colors from "tailwindcss/colors";
import { fetchAppPostRideEvaluation } from "~/apiClient/apiComponents";
import { ToIcon } from "~/components/icon/to";
import { Button } from "~/components/primitives/button/button";
import { Rating } from "~/components/primitives/rating/rating";
import { Text } from "~/components/primitives/text/text";
import { useClientAppRequestContext } from "~/contexts/user-context";
import { isClientApiError } from "~/types";

export const Arrived = () => {
  const { payload } = useClientAppRequestContext();
  const [rating, setRating] = useState(0);
  const navigate = useNavigate();
  const modalRef = useRef<{ close: () => void }>(null);

  const onClick: MouseEventHandler<HTMLButtonElement> = useCallback(
    (e) => {
      e.preventDefault();
      try {
        void (async () => {
          try {
            await fetchAppPostRideEvaluation({
              pathParams: {
                rideId: payload?.ride_id ?? "",
              },
              body: {
                evaluation: rating,
              },
            });
          } catch (error) {
            if (isClientApiError(error)) {
              if (error.stack.status === 400)
                [navigate("/client/register-payment")];
            }
          } finally {
            if (modalRef.current) {
              modalRef.current.close();
            }
          }
        })();
      } catch (e) {
        if (isClientApiError(e)) {
          console.error(`CONSOLE ERROR: ${e.message}`);
        }
      }
    },
    [payload, rating, modalRef, navigate],
  );

  return (
    <Form className="h-full flex flex-col items-center justify-center">
      <div className="flex flex-col items-center gap-6 mb-14">
        <ToIcon className="size-[90px]" color={colors.red[500]} />
        <Text size="xl">目的地に到着しました</Text>
      </div>
      <div className="flex flex-col items-center w-80">
        <Text className="mb-4">今回のドライブはいかがでしたか？</Text>
        <Rating
          name="rating"
          rating={rating}
          setRating={setRating}
          className="mb-10"
        />
        <Button
          variant="primary"
          type="submit"
          onClick={onClick}
          className="w-full mt-1"
        >
          評価してドライビングを完了
        </Button>
      </div>
    </Form>
  );
};
