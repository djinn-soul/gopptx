from __future__ import annotations  # noqa: D100

from gopptx import Presentation, ops


def test_execute_batch_mixed_results_preserve_order_and_error_details() -> None:  # noqa: D103
    with Presentation.new("Batch Mixed") as prs:
        prs.add_slide("A")

        results = prs.execute_batch([
            {"op": ops.OP_SLIDE_COUNT, "payload": {}},
            {"op": "missing_op", "payload": {}},
            {
                "op": ops.OP_SET_SLIDE_TITLE,
                "payload": {"slide_index": 0, "title": "B"},
            },
        ])

        assert len(results) == 3  # noqa: PLR2004, S101
        assert results[0].get("ok") is True  # noqa: S101
        assert results[0].get("op") == ops.OP_SLIDE_COUNT  # noqa: S101

        assert results[1].get("ok") is False  # noqa: S101
        assert results[1].get("op") == "missing_op"  # noqa: S101
        err = results[1].get("error", {})
        assert err.get("code") == "UNKNOWN_OP"  # noqa: S101
        assert err.get("details", {}).get("index") == 1  # noqa: S101

        assert results[2].get("ok") is True  # noqa: S101
        assert results[2].get("op") == ops.OP_SET_SLIDE_TITLE  # noqa: S101
        assert prs.slides[0].title == "B"  # noqa: S101


def test_execute_batch_stop_on_error_stops_following_commands() -> None:  # noqa: D103
    with Presentation.new("Batch Stop") as prs:
        prs.add_slide("A")
        initial_title = prs.slides[0].title

        results = prs.execute_batch(
            [
                {"op": "missing_op", "payload": {}},
                {
                    "op": ops.OP_SET_SLIDE_TITLE,
                    "payload": {"slide_index": 0, "title": "B"},
                },
            ],
            stop_on_error=True,
        )

        assert len(results) == 1  # noqa: S101
        assert results[0].get("ok") is False  # noqa: S101
        assert prs.slides[0].title == initial_title  # noqa: S101
