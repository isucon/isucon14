import { Form } from "@remix-run/react";
import { MouseEventHandler, useCallback, useEffect, useState } from "react";
import colors from "tailwindcss/colors";
import { fetchAppPostRideEvaluation } from "~/api/api-components";
import { PinIcon } from "~/components/icon/pin";
import { Price } from "~/components/modules/price/price";
import { Button } from "~/components/primitives/button/button";
import { ClickableRating } from "~/components/primitives/rating/clickable-rating";
import { Text } from "~/components/primitives/text/text";
import { useClientContext } from "~/contexts/client-context";

import confetti from "canvas-confetti";

export const Arrived = ({ onEvaluated }: { onEvaluated: () => void }) => {
  const { data } = useClientContext();
  const [rating, setRating] = useState(0);
  const [errorMessage, setErrorMessage] = useState<string>();

  const onClick: MouseEventHandler<HTMLButtonElement> = useCallback(
    (e) => {
      e.preventDefault();
      if (rating < 1 || rating > 5) {
        setErrorMessage("評価は1から5の間でなければなりません。");
        return;
      }
      try {
        void fetchAppPostRideEvaluation({
          pathParams: {
            rideId: data?.ride_id ?? "",
          },
          body: {
            evaluation: rating,
          },
        });
      } catch (error) {
        console.error(error);
      }
      onEvaluated();
    },
    [onEvaluated, data?.ride_id, rating],
  );

  useEffect(() => {
    void confetti({
      origin: { y: 0.7 },
      spread: 60,
      colors: [
        colors.yellow[500],
        colors.cyan[300],
        colors.green[500],
        colors.indigo[500],
        colors.red[500],
      ],
    });
  }, []);

  return (
    <>
      <Form className="w-full h-full flex flex-col items-center justify-center max-w-md mx-auto">
        <div className="flex flex-col items-center gap-6 mb-14">
          <PinIcon className="size-[90px]" color={colors.red[500]} />
          <Text size="xl" bold>
            目的地に到着しました
          </Text>
        </div>
        <div className="flex flex-col items-center w-80">
          <Text className="mb-4">今回のドライブはいかがでしたか？</Text>
          <ClickableRating
            name="rating"
            rating={rating}
            setRating={setRating}
            className="mb-10"
          />
          {data?.fare && (
            <Price pre="運賃" value={data.fare} className="mb-6"></Price>
          )}
          <Button
            variant="primary"
            type="submit"
            onClick={onClick}
            className="w-full"
          >
            評価して料金を支払う
          </Button>
        </div>
        {errorMessage && (
          <Text variant="danger" className="mt-2">
            {errorMessage}
          </Text>
        )}
      </Form>
    </>
  );
};
