from __future__ import annotations

from gopptx import Presentation, ops


def test_execute_batch_mixed_results_preserve_order_and_error_details():
    with Presentation.new("Batch Mixed") as prs:
        prs.add_slide("A")

        results = prs.execute_batch(
            [
                {"op": ops.OP_SLIDE_COUNT, "payload": {}},
                {"op": "missing_op", "payload": {}},
                {
                    "op": ops.OP_SET_SLIDE_TITLE,
                    "payload": {"slide_index": 0, "title": "B"},
                },
            ]
        )

        assert len(results) == 3
        assert results[0].get("ok") is True
        assert results[0].get("op") == ops.OP_SLIDE_COUNT

        assert results[1].get("ok") is False
        assert results[1].get("op") == "missing_op"
        err = results[1].get("error", {})
        assert err.get("code") == "UNKNOWN_OP"
        assert err.get("details", {}).get("index") == 1

        assert results[2].get("ok") is True
        assert results[2].get("op") == ops.OP_SET_SLIDE_TITLE
        assert prs.slides[0].title == "B"


def test_execute_batch_stop_on_error_stops_following_commands():
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

        assert len(results) == 1
        assert results[0].get("ok") is False
        assert prs.slides[0].title == initial_title
